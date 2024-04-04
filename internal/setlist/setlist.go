package setlist

import (
	"bytes"
	"log"

	"github.com/laenzlinger/setlist/internal/song"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

type Section struct {
	Header string
	Songs  []song.Song
}

type Setlist struct {
	Sections    []Section
	Source      []byte
	TableHeader ast.Node
	Markdown    goldmark.Markdown
}

func (sl Setlist) Render() string {
	doc := sl.generate()
	var buf bytes.Buffer

	err := sl.Markdown.Renderer().Render(&buf, sl.Source, doc)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (sl Setlist) generate() *ast.Document {
	doc := ast.NewDocument()
	table := east.NewTable()
	doc.AppendChild(doc, table)
	if sl.TableHeader != nil {
		table.AppendChild(table, sl.TableHeader)
	}

	for _, section := range sl.Sections {
		if sl.TableHeader == nil {
			p := ast.NewParagraph()
			p.AppendChild(p, ast.NewString([]byte(section.Header)))
			doc.AppendChild(doc, p)
		}
		for _, song := range section.Songs {
			table.AppendChild(table, song.TableRow)
		}
		if sl.TableHeader == nil {
			doc.AppendChild(doc, table)
			table = east.NewTable()
		}
	}
	return doc
}
