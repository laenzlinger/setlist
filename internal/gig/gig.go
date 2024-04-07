package gig

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/laenzlinger/setlist/internal/config"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Section struct {
	Header     []byte
	SongTitles []string
}

type Gig struct {
	Name     string
	Sections []Section
}

func New(band config.Band, gig string) (Gig, error) {
	file, err := os.Open(filepath.Join(band.Source, "Gigs", gig+".md"))
	if err != nil {
		return Gig{}, fmt.Errorf("failed to open Gig: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	content, err := io.ReadAll(file)
	if err != nil {
		return Gig{}, fmt.Errorf("failed to read Gig: %w", err)
	}

	gigName := fmt.Sprintf("%s @ %s", band.Name, gig)
	return parse(gigName, content), nil
}

func parse(gigName string, content []byte) Gig {
	markdown := goldmark.New()
	doc := markdown.Parser().Parse(text.NewReader(content))
	result := Gig{
		Name:     gigName,
		Sections: []Section{{}},
	}
	i := 0
	headerStart := 0
	headerStop := 0
	for first := doc.FirstChild(); first != nil; first = first.NextSibling() {
		if first.Kind() == ast.KindList {
			result.Sections[i].Header = content[headerStart:headerStop]
			for second := first.FirstChild(); second != nil; second = second.NextSibling() {
				t := string(second.Text(content))
				result.Sections[i].SongTitles = append(result.Sections[i].SongTitles, t)
			}
			_ = ast.Walk(first, func(n ast.Node, _ bool) (ast.WalkStatus, error) {
				if n.Type() == ast.TypeInline {
					return ast.WalkContinue, nil
				}
				headerStart = segmentStop(n, headerStart)
				return ast.WalkContinue, nil
			})
		} else {
			headerStop = segmentStop(first, headerStop)
		}
		if len(result.Sections[i].SongTitles) > 0 && first.NextSibling() != nil {
			i++
			result.Sections = append(result.Sections, Section{})
		}
	}
	return result
}

func segmentStop(n ast.Node, fallback int) int {
	segLen := n.Lines().Len()
	if segLen > 0 {
		return n.Lines().At(segLen - 1).Stop
	}
	return fallback
}
