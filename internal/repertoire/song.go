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

func (s Song) RemoveRows(indexes map[int]bool) Song {
	i := 0
	toRemove := []ast.Node{}
	for r := s.TableRow.FirstChild(); r != nil; r = r.NextSibling() {
		if indexes[i] {
			toRemove = append(toRemove, r)
		}
		i++
	}
	for _, r := range toRemove {
		s.TableRow.RemoveChild(s.TableRow, r)
	}
	return s
}
