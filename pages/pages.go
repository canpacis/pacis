package pages

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	h "github.com/canpacis/pacis/ui/html"
)

type ctxkey string

func Set(ctx context.Context, key string, value any) context.Context {
	return context.WithValue(ctx, ctxkey(fmt.Sprintf("%s:%s", "app", key)), value)
}

func Get[T any](ctx context.Context, key string) T {
	value := ctx.Value(ctxkey(fmt.Sprintf("%s:%s", "app", key)))
	cast, ok := value.(T)
	if !ok {
		var v T
		log.Fatalf("failed to cast ctx key '%s' to %T\n", key, v)
		return v
	}
	return cast
}

type PageContext struct {
	w http.ResponseWriter
	r *http.Request
}

func (ctx *PageContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.r.Context().Deadline()
}

func (ctx *PageContext) Err() error {
	return ctx.r.Context().Err()
}
func (ctx *PageContext) Value(key any) any {
	return ctx.r.Context().Value(key)
}

func (ctx *PageContext) Done() <-chan struct{} {
	return ctx.r.Context().Done()
}

func (ctx *PageContext) Request() *http.Request {
	return ctx.r
}

func (ctx *PageContext) Redirect(to string) h.I {
	http.Redirect(ctx.w, ctx.r, to, http.StatusFound)
	return h.Frag()
}

func (ctx *PageContext) NotFound() h.I {
	ctx.w.Header().Set("Content-Type", "text/html")
	ctx.w.WriteHeader(http.StatusNotFound)

	return NotFoundPage(ctx)
}

func (ctx *PageContext) Set(key string, value any) {
	c := context.WithValue(ctx, ctxkey(fmt.Sprintf("%s:%s", "app", key)), value)
	ctx.r = ctx.r.Clone(c)
}

type Page func(*PageContext) h.I

type LayoutContext struct {
	*PageContext
	head   h.I
	outlet h.I
}

func (ctx LayoutContext) Head() h.I {
	return ctx.head
}

func (ctx LayoutContext) Outlet() h.I {
	return ctx.outlet
}

type Layout func(*LayoutContext) h.I

func WrapLayout(layout Layout, rest ...Layout) Layout {
	switch len(rest) {
	case 0:
		return layout
	case 1:
		return func(lc *LayoutContext) h.I {
			lc.outlet = layout(lc)
			return rest[0](lc)
		}
	default:
		first := func(lc *LayoutContext) h.I {
			lc.outlet = layout(lc)
			return rest[0](lc)
		}
		return WrapLayout(first, rest[1:]...)
	}
}

var NotFoundPage Page = func(pc *PageContext) h.I {
	return h.P(h.Text("Not Found"))
}

func SetNotFoundPage(page Page) {
	NotFoundPage = page
}

var assetmap = map[string]string{}

func RegisterAssets(assets map[string]string) {
	assetmap = assets
}

func Asset(src string) string {
	resolved, ok := assetmap[src]
	if !ok {
		log.Fatalf("failed to resolve asset %s", src)
	}
	return resolved
}

type Route interface {
	http.Handler
	Path() string
}

type HomeRoute struct {
	page        Page
	layout      Layout
	head        h.I
	middlewares []func(http.Handler) http.Handler
}

func (HomeRoute) Path() string {
	return "/"
}

func (hr *HomeRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		ctx := &PageContext{w: w, r: r}
		var renderer h.I
		var page Page

		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			page = NotFoundPage
		} else {
			w.WriteHeader(http.StatusOK)
			page = hr.page
		}

		if hr.layout != nil {
			renderer = hr.layout(&LayoutContext{PageContext: ctx, head: hr.head, outlet: page(ctx)})
		} else {
			renderer = page(ctx)
		}
		renderer.Render(ctx, w)
	})
	for _, middleware := range hr.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
}

func NewHomeRoute(page Page, layout Layout, head h.I, middlewares ...func(http.Handler) http.Handler) *HomeRoute {
	return &HomeRoute{page: page, layout: layout, head: head, middlewares: middlewares}
}

type PageRoute struct {
	path        string
	page        Page
	layout      Layout
	head        h.I
	middlewares []func(http.Handler) http.Handler
}

func NewPageRoute(path string, page Page, layout Layout, head h.I, middlewares ...func(http.Handler) http.Handler) *PageRoute {
	return &PageRoute{path: path, page: page, layout: layout, head: head, middlewares: middlewares}
}

func (pr PageRoute) Path() string {
	return pr.path
}

func (pr *PageRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		ctx := &PageContext{w: w, r: r}
		var renderer h.I
		if pr.layout != nil {
			renderer = pr.layout(&LayoutContext{PageContext: ctx, head: pr.head, outlet: pr.page(ctx)})
		} else {
			renderer = pr.page(ctx)
		}
		renderer.Render(ctx, w)
	})
	for _, middleware := range pr.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
}

type RedirectRoute struct {
	path string
	To   string
	Code int
}

func NewRedirectRoute(path, to string, code int) *RedirectRoute {
	return &RedirectRoute{path: path, To: to, Code: code}
}

func (rr RedirectRoute) Path() string {
	return rr.path
}

func (rr *RedirectRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, rr.To, rr.Code)
}

type RawRoute struct {
	path        string
	contenttyp  string
	content     []byte
	middlewares []func(http.Handler) http.Handler
}

func NewRawRoute(path, typ string, content []byte, middlewares ...func(http.Handler) http.Handler) *RawRoute {
	return &RawRoute{path: path, contenttyp: typ, content: content, middlewares: middlewares}
}

func (rr RawRoute) Path() string {
	return rr.path
}

func (rr *RawRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", rr.contenttyp)
		w.WriteHeader(http.StatusOK)
		w.Write(rr.content)
	})
	for _, middleware := range rr.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
}
