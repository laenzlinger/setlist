package song

import (
	"github.com/laenzlinger/setlist/internal/nodetext"
	"github.com/yuin/goldmark/ast"
)

type Song struct {
	TableRow ast.Node
	Title    string
}

func New(ast ast.Node, source []byte) Song {
	col := ast.FirstChild()
	return Song{TableRow: ast, Title: nodetext.Extract(col, source)}
}

func (s Song) String() string {
	return s.Title
}

func (s Song) RemoveColumns(idx Indexes) Song {
	s.TableRow = RemoveCols(idx, s.TableRow)
	return s
}
