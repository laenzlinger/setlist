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
	band        string
	song        string
	placeholder bool
}

func AllForBand(band string) error {
	songs := map[string]bool{}
	sheets := Sheet{band: band}
	files, err := os.ReadDir(sheets.sourceDir())
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
		return fmt.Errorf("no songs found in %s", sheets.sourceDir())
	}
	songNames := []string{}
	for song := range songs {
		songNames = append(songNames, song)
	}
	sort.Strings(songNames)
	return forSongs(band, songNames, fmt.Sprintf("All %s Songs", band))
}

func ForGig(band string, gig gig.Gig) error {
	songs := []string{}
	for _, section := range gig.Sections {
		songs = append(songs, section.SongTitles...)
	}
	return forSongs(band, songs, gig.Name)
}

func forSongs(band string, songs []string, sheetName string) error {
	files := []string{}
	for _, song := range songs {
		s := &Sheet{band: band, song: song}
		err := s.verifySheetPdf()
		if err != nil {
			return fmt.Errorf("failed to create sheet PDF for `%s - %s : %w", band, song, err)
		}
		files = append(files, s.pdfName())
	}

	tmpl.PrepareOut()
	target := fmt.Sprintf("out/Cheat Sheet %v.pdf", sheetName)

	err := pdf.MergeCreateFile(files, target, false, nil)
	if err != nil {
		return fmt.Errorf("failed to merge PDF files: %w", err)
	}
	return nil
}

func (s *Sheet) verifySheetPdf() error {
	sourceExists, targetExists := true, true

	source, err := os.Stat(s.sourceName())
	if errors.Is(err, os.ErrNotExist) {
		sourceExists = false
	} else if err != nil {
		return err
	}

	target, err := os.Stat(s.pdfName())
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
	log.Printf("generate from source for `%s`", s.song)
	buf := bytes.NewBuffer([]byte{})
	args := []string{"--headless", "--convert-to", "pdf", "--outdir", s.sourceDir(), s.sourceName()}
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

func (s *Sheet) pdfName() string {
	return fmt.Sprintf("%s/%s.pdf", s.pdfDir(), s.song)
}

func (s *Sheet) sourceName() string {
	return fmt.Sprintf("%s/%s.odt", s.sourceDir(), s.song)
}

func (s *Sheet) pdfDir() string {
	if s.placeholder {
		return "out/placeholder"
	}
	return s.sourceDir()
}

func (s *Sheet) sourceDir() string {
	return fmt.Sprintf("%s/Songs", s.band)
}

func (s *Sheet) generatePlaceholder() error {
	s.placeholder = true
	if err := os.MkdirAll(s.pdfDir(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create out directory: %w", err)
	}
	//nolint: gosec // content does not contain html
	filename, err := tmpl.CreatePlaceholder(&tmpl.Data{Content: template.HTML(s.song), Title: s.song})
	if err != nil {
		return err
	}

	return convert.HTMLToPDF(filename, s.pdfName())
}
