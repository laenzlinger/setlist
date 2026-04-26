package nodetext_test

import (
	"testing"

	"github.com/laenzlinger/setlist/internal/nodetext"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func firstChild(source []byte, kind ast.NodeKind) ast.Node {
	doc := goldmark.New().Parser().Parse(text.NewReader(source))
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if n.Kind() == kind {
			return n
		}
	}
	return nil
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
		kind   ast.NodeKind
	}{
		{
			name:   "heading",
			source: "# Hello World",
			kind:   ast.KindHeading,
			want:   "Hello World",
		},
		{
			name:   "paragraph",
			source: "some text",
			kind:   ast.KindParagraph,
			want:   "some text",
		},
		{
			name:   "list item",
			source: "* item one",
			kind:   ast.KindList,
			want:   "item one",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.source)
			n := firstChild(source, tt.kind)
			assert.NotNil(t, n)
			got := nodetext.Extract(n, source)
			assert.Equal(t, tt.want, got)
		})
	}
}
