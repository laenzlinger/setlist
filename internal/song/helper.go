package song

import "github.com/yuin/goldmark/ast"

type Indexes map[int]bool

func RemoveCols(idx Indexes, row ast.Node) ast.Node {
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
