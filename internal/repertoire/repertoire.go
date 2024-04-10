package repertoire

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/laenzlinger/setlist/internal/config"
	"github.com/laenzlinger/setlist/internal/gig"
	"github.com/laenzlinger/setlist/internal/setlist"
	"github.com/laenzlinger/setlist/internal/song"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/wikilink"
)

type Repertoire struct {
	songs    []song.Song
	header   song.Header
	columns  []string
	source   []byte
	markdown goldmark.Markdown
}

func New(band config.Band) (Repertoire, error) {
	file, err := os.Open(path.Join(band.Source, "Repertoire.md"))
	if err != nil {
		return Repertoire{}, fmt.Errorf("failed to open Repertoire: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return Repertoire{}, fmt.Errorf("failed to read Repertoire: %w", err)
	}

	return from(content), nil
}

func from(source []byte) Repertoire {
	md := goldmark.New(goldmark.WithExtensions(
		extension.GFM,
		&wikilink.Extender{},
	))
	doc := md.Parser().Parse(text.NewReader(source))
	result := Repertoire{source: source, markdown: md}

	table := doc.FirstChild()
	for row := table.FirstChild(); row != nil; row = row.NextSibling() {
		if row.Kind() == east.KindTableRow {
			result.songs = append(result.songs, song.New(row, source))
		}
		if row.Kind() == east.KindTableHeader {
			result.header = song.NewHeader(&row)
			for h := row.FirstChild(); h != nil; h = h.NextSibling() {
				result.columns = append(result.columns, string(h.Text(source)))
			}
		}
	}

	return result
}

func (rep Repertoire) Merge(g gig.Gig) setlist.Setlist {
	sections := []setlist.Section{}
	for _, section := range g.Sections {
		sect := setlist.Section{Header: section.HeaderText()}
		for _, title := range section.SongTitles {
			found := false
			for _, song := range rep.songs {
				if normalize(song.Title) == normalize(title) {
					sect.Songs = append(sect.Songs, song)
					found = true
				}
			}
			if !found {
				log.Fatalf("Song `%s` not found in repertoire", title)
			}
		}
		sections = append(sections, sect)
	}

	return setlist.Setlist{
		Sections:    sections,
		Source:      rep.source,
		TableHeader: rep.header,
		Markdown:    rep.markdown,
	}
}

func (rep Repertoire) ExcludeColumns(columns ...string) Repertoire {
	idx := song.Indexes{}
	for _, toRemove := range columns {
		for i, c := range rep.columns {
			if normalize(c) == normalize(toRemove) {
				idx[i] = true
			}
		}
	}
	for _, song := range rep.songs {
		song.RemoveColumns(idx)
	}
	rep.header = rep.header.RemoveColumns(idx)
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

func (rep Repertoire) NoHeader() Repertoire {
	rep.header = rep.header.Remove()
	return rep
}

var valid = regexp.MustCompile(`[^a-z]+`)

func normalize(n string) string {
	return valid.ReplaceAllString(strings.ToLower(n), "")
}
