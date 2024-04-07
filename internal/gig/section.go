package gig

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Section struct {
	Header     []byte
	SongTitles []string
}

func NewSection() Section {
	return Section{}
}

func (s Section) HeaderHTML() (string, error) {
	b := bytes.NewBuffer([]byte{})
	markdown := goldmark.New()
	err := markdown.Convert(s.Header, b)
	if err != nil {
		return "", fmt.Errorf("failed to convert header to HTML: %w", err)
	}

	return b.String(), nil
}

// Returns the first header text.
func (s Section) HeaderText() string {
	result := ""
	markdown := goldmark.New()
	doc := markdown.Parser().Parse(text.NewReader(s.Header))
	_ = ast.Walk(doc, func(n ast.Node, _ bool) (ast.WalkStatus, error) {
		if n.Kind() == ast.KindHeading {
			result = string(n.Text(s.Header))
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})

	return result
}
