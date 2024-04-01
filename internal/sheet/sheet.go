package sheet

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
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

func ForGig(band config.Band, gig gig.Gig) error {
	sheets := []Sheet{}
	for i, section := range gig.Sections {
		header := Sheet{band: band, name: fmt.Sprintf("section-header-%d", i), content: section.Header}
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
		err := s.createPdf()
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

	return os.RemoveAll(config.PlaceholderDir())
}

func (s *Sheet) createPdf() error {
	sourceExists, targetExists := true, true

	source, err := os.Stat(s.sourceFilePath())
	if errors.Is(err, os.ErrNotExist) {
		sourceExists = false
	} else if err != nil {
		return err
	}

	target, err := os.Stat(s.pdfFilePath())
	if errors.Is(err, os.ErrNotExist) {
		targetExists = false
	} else if err != nil {
		return err
	}

	if sourceExists {
		if !targetExists || target.ModTime().Before(source.ModTime()) {
			return s.generateFromSource()
		}
	}

	if !targetExists {
		return s.generatePlaceholder()
	}
	return nil
}

func (s *Sheet) generateFromSource() error {
	log.Printf("generate from source for `%s`", s.name)
	buf := bytes.NewBuffer([]byte{})
	args := []string{"--headless", "--convert-to", "pdf", "--outdir", s.sourceDir(), s.sourceFilePath()}
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

func (s *Sheet) pdfFilePath() string {
	return filepath.Join(s.pdfDir(), s.name+".pdf")
}

func (s *Sheet) sourceFilePath() string {
	return filepath.Join(s.sourceDir(), s.name+".odt")
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
	if err := os.MkdirAll(s.pdfDir(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create out directory: %w", err)
	}
	//nolint: gosec // content does not contain html
	filename, err := tmpl.CreatePlaceholder(&tmpl.Data{Content: template.HTML(s.content), Title: s.name})
	if err != nil {
		return err
	}

	return convert.HTMLToPDF(filename, s.pdfFilePath())
}
