package song

import (
	"github.com/yuin/goldmark/ast"
)

type Header struct {
	TableHeader *ast.Node
}

type Song struct {
	Title    string
	TableRow ast.Node
}

func New(ast ast.Node, source []byte) Song {
	col := ast.FirstChild()
	return Song{TableRow: ast, Title: string(col.Text(source))}
}

func (s Song) String() string {
	return s.Title
}

func (s Song) RemoveColumns(idx Indexes) Song {
	s.TableRow = RemoveCols(idx, s.TableRow)
	return s
}

func (h Header) Remove() Header {
	h.TableHeader = nil
	return h
}

func (h Header) RemoveColumns(idx Indexes) Header {
	if h.TableHeader != nil {
		removed := RemoveCols(idx, *h.TableHeader)
		h.TableHeader = &removed
	}
	return h
}

func (h Header) Empty() bool {
	return h.TableHeader == nil
}
