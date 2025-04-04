package pages

import (
	"context"
	"net/http"

	c "github.com/canpacis/pacis/ui/components"
	h "github.com/canpacis/pacis/ui/html"
)

type PageContext struct {
	context.Context
	req    *http.Request
	cancel context.CancelFunc
}

func (ctx PageContext) Request() *http.Request {
	return ctx.req
}

type Page func(*PageContext) h.I

type LayoutContext struct {
	*PageContext
	head   *c.AppHead
	outlet h.I
}

func (lc LayoutContext) Head() *c.AppHead {
	return lc.head
}

func (lc LayoutContext) Outlet() h.I {
	return lc.outlet
}

type Layout func(*LayoutContext) h.I

type Route interface {
	http.Handler
	Path() string
}

type PageRoute struct {
	path    string
	page    Page
	errview Page
	layout  Layout
}

func (pr PageRoute) Path() string {
	return pr.path
}

func (pr PageRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	pctx := &PageContext{Context: ctx, cancel: cancel, req: r}

	var renderer h.I
	if pr.layout != nil {
		head := c.CreateHead("/ui/")
		lctx := &LayoutContext{PageContext: pctx, head: head, outlet: pr.page(pctx)}
		renderer = pr.layout(lctx)
	} else {
		renderer = pr.page(pctx)
	}
	renderer.Render(pctx, w)
}

func NewPageRoute(path string, layout Layout, page Page) *PageRoute {
	return &PageRoute{path: path, page: page, layout: layout}
}
