package repertoire

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/laenzlinger/setlist/internal/gig"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

type Repertoire struct {
	songs    []Song
	columns  []string
	source   []byte
	markdown goldmark.Markdown
}

func New(band string) (Repertoire, error) {
	file, err := os.Open(fmt.Sprintf("%s/Repertoire.md", band))
	if err != nil {
		return Repertoire{}, fmt.Errorf("failed to open Repertoire file: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return Repertoire{}, fmt.Errorf("failed to read Repertoire file: %w", err)
	}

	return from(content), nil
}

func from(source []byte) Repertoire {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	doc := md.Parser().Parse(text.NewReader(source))
	result := Repertoire{source: source, markdown: md}

	table := doc.FirstChild()
	for row := table.FirstChild(); row != nil; row = row.NextSibling() {
		if row.Kind() == east.KindTableRow {
			result.songs = append(result.songs, SongFrom(row, source))
		}
		if row.Kind() == east.KindTableHeader {
			for h := row.FirstChild(); h != nil; h = h.NextSibling() {
				result.columns = append(result.columns, string(h.Text(source)))
			}
		}
	}

	return result
}

func (rep Repertoire) Filter(g gig.Gig) Repertoire {
	result := []Song{}
	for _, section := range g.Sections {
		for _, title := range section.SongTitles {
			found := false
			for _, song := range rep.songs {
				if normalize(song.Title) == normalize(title) {
					result = append(result, song)
					found = true
				}
			}
			if !found {
				log.Fatalf("Song `%s` not found in repertoire", title)
			}
		}
	}
	rep.songs = result
	return rep
}

func (rep Repertoire) Render() string {
	doc := rep.generate()
	var buf bytes.Buffer

	err := rep.markdown.Renderer().Render(&buf, rep.source, doc)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (rep Repertoire) ExcludeColumns(columns ...string) Repertoire {
	indexes := map[int]bool{}
	for _, toRemove := range columns {
		for i, c := range rep.columns {
			if normalize(c) == normalize(toRemove) {
				indexes[i] = true
			}
		}
	}
	for _, song := range rep.songs {
		song.RemoveRows(indexes)
	}
	return rep
}

func (rep Repertoire) IncludeColumns(columns ...string) Repertoire {
	exclude := []string{}
	for _, col := range rep.columns {
		found := false
		for _, inc := range columns {
			if col == inc {
				found = true
			}
		}
		if !found {
			exclude = append(exclude, col)
		}
	}
	return rep.ExcludeColumns(exclude...)
}

func (rep Repertoire) generate() *ast.Document {
	doc := ast.NewDocument()
	table := east.NewTable()
	doc.AppendChild(doc, table)
	for _, song := range rep.songs {
		table.AppendChild(table, song.TableRow)
	}
	return doc
}

var valid = regexp.MustCompile(`[^a-z]+`)

func normalize(n string) string {
	return valid.ReplaceAllString(strings.ToLower(n), "")
}
