package repertoire

import (
	"github.com/yuin/goldmark/ast"
)

type Song struct {
	Title    string
	TableRow ast.Node
}

func SongFrom(ast ast.Node, source []byte) Song {
	col := ast.FirstChild()
	return Song{TableRow: ast, Title: string(col.Text(source))}
}

func (s Song) String() string {
	return s.Title
}

func (s Song) removeColumns(idx indexes) Song {
	s.TableRow = removeCols(idx, s.TableRow)
	return s
}
