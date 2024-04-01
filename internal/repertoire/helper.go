package repertoire

import "github.com/yuin/goldmark/ast"

type indexes map[int]bool

func removeCols(idx indexes, row ast.Node) ast.Node {
	i := 0
	toRemove := []ast.Node{}
	for r := row.FirstChild(); r != nil; r = r.NextSibling() {
		if idx[i] {
			toRemove = append(toRemove, r)
		}
		i++
	}
	for _, r := range toRemove {
		row.RemoveChild(row, r)
	}
	return row
}
