package sheet

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/laenzlinger/setlist/internal/config"
	"github.com/laenzlinger/setlist/internal/gig"
	convert "github.com/laenzlinger/setlist/internal/html/pdf"
	tmpl "github.com/laenzlinger/setlist/internal/html/template"
	pdf "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Sheet struct {
	band        config.Band
	name        string
	content     string
	placeholder bool
}

func AllForBand(band config.Band) error {
	songs := map[string]bool{}
	aSheet := Sheet{band: band}
	files, err := os.ReadDir(aSheet.sourceDir())
	if err != nil {
		return fmt.Errorf("failed to list Band directory: %w", err)
	}

	for _, file := range files {
		extraw := filepath.Ext(file.Name())
		ext := strings.ToLower(extraw)
		if !file.IsDir() && (ext == ".pdf" || ext == ".odt") {
			songs[strings.TrimSuffix(filepath.Base(file.Name()), ext)] = true
		}
	}
	if len(songs) == 0 {
		return fmt.Errorf("no songs found in %s", aSheet.sourceDir())
	}
	songNames := []string{}
	for song := range songs {
		songNames = append(songNames, song)
	}
	sort.Strings(songNames)
	sheets := []Sheet{}
	for _, title := range songNames {
		s := Sheet{band: band, name: title, content: title}
		sheets = append(sheets, s)
	}
	return merge(sheets, fmt.Sprintf("for all %s songs", band.Name))
}

const SectionPrefix = "SECTION:"

type sectionHeaders map[string]int

func (sh sectionHeaders) add(value string) {
	sh[value]++
}

func (sh sectionHeaders) filename(value string) string {
	if sh[value] <= 1 {
		return SectionPrefix + value
	}
	return fmt.Sprintf("%s%s %d", SectionPrefix, value, sh[value])
}

func ForGig(band config.Band, gig gig.Gig) error {
	if err := os.MkdirAll(config.PlaceholderDir(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(config.PlaceholderDir())

	sheets := []Sheet{}
	sh := sectionHeaders{}
	for _, section := range gig.Sections {
		h := section.HeaderText()
		sh.add(h)
		html, err := section.HeaderHTML()
		if err != nil {
			return err
		}
		header := Sheet{band: band, name: sh.filename(h), content: html}
		sheets = append(sheets, header)
		for _, title := range section.SongTitles {
			song := Sheet{band: band, name: title, content: title}
			sheets = append(sheets, song)
		}
	}
	return merge(sheets, gig.Name)
}

func merge(sheets []Sheet, outputFileName string) error {
	files := []string{}
	for _, s := range sheets {
		err := s.ensurePdf()
		if err != nil {
			return fmt.Errorf("failed to create sheet PDF for `%s`: %w", s.name, err)
		}
		files = append(files, s.pdfFilePath())
	}

	tmpl.PrepareTarget()
	target := filepath.Join(config.Target(), fmt.Sprintf("Cheat Sheet %v.pdf", outputFileName))

	err := pdf.MergeCreateFile(files, target, false, nil)
	if err != nil {
		return fmt.Errorf("failed to merge PDF files: %w", err)
	}

	err = cleanupBookmarks(target)
	if err != nil {
		return fmt.Errorf("failed cleanup bookmarks: %w", err)
	}

	return nil
}

// Create or update the pdf from the source.
func (s *Sheet) ensurePdf() error {
	mdExists, odtExists, targetExists := true, true, true

	odt, err := os.Stat(s.odtFilePath())
	if errors.Is(err, os.ErrNotExist) {
		odtExists = false
	} else if err != nil {
		return err
	}

	md, err := os.Stat(s.mdFilePath())
	if errors.Is(err, os.ErrNotExist) {
		mdExists = false
	} else if err != nil {
		return err
	}

	target, err := os.Stat(s.pdfFilePath())
	if errors.Is(err, os.ErrNotExist) {
		targetExists = false
	} else if err != nil {
		return err
	}

	if odtExists {
		if !targetExists || target.ModTime().Before(odt.ModTime()) {
			return s.generateFromOdt()
		}
	} else if mdExists {
		if !targetExists || target.ModTime().Before(md.ModTime()) {
			return s.generateFromMarkdown()
		}
	}

	if !targetExists {
		return s.generatePlaceholder()
	}
	return nil
}

func (s *Sheet) generateFromOdt() error {
	log.Printf("generate from odt source for `%s`", s.name)
	buf := bytes.NewBuffer([]byte{})
	args := []string{"--headless", "--convert-to", "pdf", "--outdir", s.sourceDir(), s.odtFilePath()}
	if config.RunningInContainer() {
		args = append(args, fmt.Sprintf("-env:UserInstallation=file:///%s", config.UserHome()))
	}

	cmd := exec.Command("libreoffice", args...)
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		log.Println(buf.String())
		return err
	}
	return nil
}

func (s *Sheet) generateFromMarkdown() error {
	log.Printf("generate from markdown source for `%s`", s.name)

	file, err := os.Open(s.mdFilePath())
	if err != nil {
		return fmt.Errorf("failed to open source markdown: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read Gig: %w", err)
	}

	md := goldmark.New(goldmark.WithExtensions(
		extension.GFM,
	))

	var buf bytes.Buffer
	if err = md.Convert(content, &buf); err != nil {
		return err
	}

	filename, err := tmpl.CreateSongsheet(&tmpl.Data{
		Content: template.HTML(buf.String()), //nolint: gosec // not a web application
		Title:   s.name,
	})
	if err != nil {
		return err
	}

	return convert.HTMLToPDF(filename, s.pdfFilePath())
}

func (s *Sheet) pdfFilePath() string {
	return filepath.Join(s.pdfDir(), s.name+".pdf")
}

func (s *Sheet) odtFilePath() string {
	return filepath.Join(s.sourceDir(), s.name+".odt")
}

func (s *Sheet) mdFilePath() string {
	return filepath.Join(s.sourceDir(), s.name+".md")
}

func (s *Sheet) pdfDir() string {
	if s.placeholder {
		return config.PlaceholderDir()
	}
	return s.sourceDir()
}

func (s *Sheet) sourceDir() string {
	return filepath.Join(s.band.Source, "Songs")
}

func (s *Sheet) generatePlaceholder() error {
	s.placeholder = true
	//nolint: gosec // content does not contain html
	filename, err := tmpl.CreatePlaceholder(&tmpl.Data{Content: template.HTML(s.content), Title: s.name})
	if err != nil {
		return err
	}

	return convert.HTMLToPDF(filename, s.pdfFilePath())
}

func cleanupBookmarks(source string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}

	bms, err := pdf.Bookmarks(in, nil)
	if err != nil {
		return err
	}

	partitioned := false
	newBms := []pdfcpu.Bookmark{}
	var currentSection *pdfcpu.Bookmark
	for i := range bms {
		sectionStart := strings.HasPrefix(bms[i].Title, SectionPrefix)
		bms[i].Title = strings.TrimPrefix(strings.TrimSuffix(bms[i].Title, ".pdf"), SectionPrefix)

		if sectionStart {
			partitioned = true
		}
		switch {
		case partitioned && sectionStart:
			if currentSection != nil {
				newBms = append(newBms, *currentSection)
			}
			currentSection = &bms[i]
		case partitioned && !sectionStart:
			currentSection.Kids = append(currentSection.Kids, bms[i])
		default:
			newBms = append(newBms, bms[i])
		}
	}
	if partitioned && currentSection != nil {
		newBms = append(newBms, *currentSection)
	}

	err = pdf.AddBookmarksFile(source, source, newBms, true, nil)
	if err != nil {
		return err
	}

	return nil
}
