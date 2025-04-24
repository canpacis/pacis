package async

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"

	c "github.com/canpacis/pacis/ui/components"
	h "github.com/canpacis/pacis/ui/html"
)

func randid() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return "pacis-" + hex.EncodeToString(buf)
}

type StreamElement struct {
	fn       func() h.Element
	fallback h.Node
}

func (n *StreamElement) Render(ctx context.Context, w io.Writer) error {
	// TODO
	// pctx, ok := ctx.(*pages.Context)
	// // If context is not a page context, render sync
	// if !ok {
	// 	return n.fn().Render(ctx, w)
	// }

	id := randid()
	// dequeue := pctx.QueueElement()

	go func(fn func() h.Element, id string) {
		// for !pctx.Ready() {
		// 	time.Sleep(time.Microsecond * 100)
		// }
		element := fn()
		element.AddAttribute(h.SlotAttr(id))
		element.AddAttribute(c.X("show", "false"))
		// dequeue(element)
	}(n.fn, id)

	placholder := h.Slot(h.Name(id))

	if n.fallback != nil {
		placholder.AddNode(n.fallback)
	}
	return placholder.Render(ctx, w)
}

func (*StreamElement) NodeType() h.NodeType {
	return h.NodeFragment
}

func Element(fn func() h.Element, fallback ...h.Node) *StreamElement {
	el := &StreamElement{fn: fn}
	if len(fallback) != 0 {
		el.fallback = fallback[0]
	}
	return el
}
