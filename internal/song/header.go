package song

import (
	"github.com/yuin/goldmark/ast"
)

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
