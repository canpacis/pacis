package server

import (
	"context"
	"io"
	"log"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/internal/util"
)

type AsyncNode struct {
	fn       func() *html.Element
	fallback html.Node
}

func (*AsyncNode) Item() {}

func (n *AsyncNode) Render(ctx context.Context, w io.Writer) error {
	context, ok := ctx.(*Context)
	if ok {
		id := util.PrefixedID("pacis")
		context.RegisterChunk(id, n.fn)
		return html.Slot(html.Name(id), html.Fragment(n.fallback)).Render(ctx, w)
	}

	return n.fn().Render(ctx, w)
}

func Async(fn func() *html.Element, fallback html.Node) *AsyncNode {
	if fallback == nil {
		fallback = html.Fragment()
	}
	return &AsyncNode{fn: fn, fallback: fallback}
}

// TODO: Implement eagerness rules
type SpeculationHook struct {
	render bool
	ctx    *Context
}

func (*SpeculationHook) Item() {}

func (h *SpeculationHook) Done(el *html.Element) {
	href, ok := el.Attributes["href"]
	if !ok {
		log.Fatal("Speculation hook used in an element without an href attribute")
	}
	if h.render {
		h.ctx.RegisterPrerender(href)
	} else {
		h.ctx.RegisterPrefetch(href)
	}
}

func Prefetch(ctx context.Context) *SpeculationHook {
	context, ok := ctx.(*Context)
	if !ok {
		log.Fatal("Speculation hook used outside of server rendering context")
	}
	return &SpeculationHook{ctx: context}
}

func Prerender(ctx context.Context) *SpeculationHook {
	context, ok := ctx.(*Context)
	if !ok {
		log.Fatal("Speculation hook used outside of server rendering context")
	}
	return &SpeculationHook{ctx: context, render: true}
}
