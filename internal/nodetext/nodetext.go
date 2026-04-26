package nodetext

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

// Extract returns the text content of a goldmark AST node
// using the non-deprecated API (ast.Text.Value instead of Node.Text).
func Extract(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			buf.Write(t.Value(source))
		} else {
			buf.WriteString(Extract(c, source))
		}
	}
	return buf.String()
}
