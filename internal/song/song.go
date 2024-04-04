package song

import (
	"github.com/yuin/goldmark/ast"
)

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

type Header struct {
	tableHeader *ast.Node
}

func NewHeader(h *ast.Node) Header {
	return Header{tableHeader: h}
}

func (h Header) Remove() Header {
	h.tableHeader = nil
	return h
}

func (h Header) RemoveColumns(idx Indexes) Header {
	if !h.Empty() {
		removed := RemoveCols(idx, *h.tableHeader)
		h.tableHeader = &removed
	}
	return h
}

func (h Header) Empty() bool {
	return h.tableHeader == nil
}

func (h Header) Node() ast.Node {
	if h.Empty() {
		return ast.NewText()
	}
	return *h.tableHeader
}
