package song_test

import (
	"testing"

	"github.com/laenzlinger/setlist/internal/song"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

func tableRow(source []byte) (*east.TableRow, []byte) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	doc := md.Parser().Parse(text.NewReader(source))
	table := doc.FirstChild()
	for r := table.FirstChild(); r != nil; r = r.NextSibling() {
		if row, ok := r.(*east.TableRow); ok {
			return row, source
		}
	}
	return nil, source
}

func TestNew(t *testing.T) {
	source := []byte("| Title | Year |\n|---|---|\n| My Song | 2024 |\n")
	row, src := tableRow(source)
	assert.NotNil(t, row)

	s := song.New(row, src)
	assert.Equal(t, "My Song", s.Title)
	assert.Equal(t, "My Song", s.String())
}

func TestRemoveCols(t *testing.T) {
	source := []byte("| A | B | C |\n|---|---|---|\n| 1 | 2 | 3 |\n")
	row, _ := tableRow(source)
	assert.NotNil(t, row)

	// Count children before.
	countBefore := 0
	for c := row.FirstChild(); c != nil; c = c.NextSibling() {
		countBefore++
	}
	assert.Equal(t, 3, countBefore)

	// Remove middle column.
	song.RemoveCols(song.Indexes{1: true}, row)

	countAfter := 0
	for c := row.FirstChild(); c != nil; c = c.NextSibling() {
		countAfter++
	}
	assert.Equal(t, 2, countAfter)
}
