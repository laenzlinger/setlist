package sheet

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"

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

func ForGig(band string, gig gig.Gig) {

	files := []string{}
	for _, section := range gig.Sections {
		for _, song := range section.SongTitles {
			s := &Sheet{band: band, song: song}
			s.verifySheetPdf()
			files = append(files, s.pdfName())
		}
	}

	tmpl.PrepareOut()
	target := fmt.Sprintf("out/Cheat Sheet %v.pdf", gig.Name)

	err := pdf.MergeCreateFile(files, target, false, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Sheet) verifySheetPdf() {
	sourceExists, targetExists := true, true

	source, err := os.Stat(s.sourceName())
	if errors.Is(err, os.ErrNotExist) {
		sourceExists = false
	} else if err != nil {
		log.Fatal(err)
	}

	target, err := os.Stat(s.pdfName())
	if errors.Is(err, os.ErrNotExist) {
		targetExists = false
	} else if err != nil {
		log.Fatal(err)
	}

	if sourceExists {
		if !targetExists || target.ModTime().Before(source.ModTime()) {
			s.generateFromSource()
			targetExists = true
		}
	}

	if !targetExists {
		s.generatePlaceholder()
	}
}

func (s *Sheet) generateFromSource() {
	log.Printf("generate from source for `%s`", s.song)
	// libreoffice --headless --convert-to pdf --outdir ${targetDir} "${sourceFile}"
	buf := bytes.NewBuffer([]byte{})
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf", "--outdir", s.sourceDir(), s.sourceName())
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		log.Println(buf.String())
		log.Fatal(err)
	}
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
	return fmt.Sprintf("%s/songs", s.band)
}

func (s *Sheet) generatePlaceholder() {
	s.placeholder = true
	if err := os.MkdirAll(s.pdfDir(), os.ModePerm); err != nil {
		log.Fatal(err)
	}
	filename := tmpl.CreatePlaceholder(&tmpl.TemplateData{Content: template.HTML(s.song), Title: s.song})
	convert.HtmlToPdf(filename, s.pdfName())
}
