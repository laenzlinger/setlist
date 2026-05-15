package setlist_test

import (
	"testing"

	"github.com/laenzlinger/setlist/internal/setlist"
	"github.com/laenzlinger/setlist/internal/song"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

func songsFromTable(source []byte) (song.Header, []song.Song) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	doc := md.Parser().Parse(text.NewReader(source))
	table := doc.FirstChild()
	var header song.Header
	var songs []song.Song
	for r := table.FirstChild(); r != nil; r = r.NextSibling() {
		if r.Kind() == east.KindTableRow {
			songs = append(songs, song.New(r, source))
		}
		if r.Kind() == east.KindTableHeader {
			header = song.NewHeader(&r)
		}
	}
	return header, songs
}

func TestRender(t *testing.T) {
	source := []byte("| Title | Year |\n|---|---|\n| Song A | 2020 |\n| Song B | 2021 |\n")
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	header, songs := songsFromTable(source)

	sl := setlist.Setlist{
		TableHeader: header,
		Markdown:    md,
		Source:      source,
		Sections: []setlist.Section{
			{Songs: songs},
		},
	}

	result := sl.Render()
	assert.Contains(t, result, "Song A")
	assert.Contains(t, result, "Song B")
	assert.Contains(t, result, "<t")
}
