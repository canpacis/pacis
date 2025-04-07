package pages

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/NYTimes/gziphandler"
	c "github.com/canpacis/pacis/ui/components"
	h "github.com/canpacis/pacis/ui/html"
)

func normalpath(p string) string {
	if len(p) == 0 || p == "/" {
		return "/"
	}
	p, _ = strings.CutPrefix(p, "/")
	p, _ = strings.CutSuffix(p, "/")
	return "/" + p
}

type ctxkey string

type Context interface {
	context.Context
	Set(string, any)
}

func Set(ctx context.Context, key string, value any) context.Context {
	return context.WithValue(ctx, ctxkey(fmt.Sprintf("%s:%s", "app", key)), value)
}

func Get[T any](ctx context.Context, key string) T {
	value := ctx.Value(ctxkey(fmt.Sprintf("%s:%s", "app", key)))
	cast, ok := value.(T)
	if !ok {
		var v T
		log.Printf("failed to cast ctx key '%s' to %T\n", key, v)
		return v
	}
	return cast
}

type PageContext struct {
	context.Context
	req    *http.Request
	cancel context.CancelFunc
}

func (ctx *PageContext) Request() *http.Request {
	return ctx.req
}

func (ctx *PageContext) Set(key string, value any) {
	c := context.WithValue(ctx, ctxkey(fmt.Sprintf("%s:%s", "app", key)), value)
	ctx.Context = c
}

type Page func(*PageContext) h.I

type LayoutContext struct {
	*PageContext
	head   *c.AppHead
	outlet h.I
}

func (ctx *LayoutContext) Set(key string, value any) {
	c := context.WithValue(ctx, ctxkey(fmt.Sprintf("%s:%s", "app", key)), value)
	ctx.Context = c
}

func (ctx LayoutContext) Head() *c.AppHead {
	return ctx.head
}

func (ctx LayoutContext) Request() *http.Request {
	return ctx.PageContext.Request()
}

func (ctx LayoutContext) Outlet() h.I {
	return ctx.outlet
}

type Layout func(*LayoutContext) h.I

type RouteItem interface {
	item()
}

type raw struct {
	bytes      []byte
	contenttyp string
}

type route struct {
	path        string
	public      *public
	root        bool
	page        Page
	layout      Layout
	redirect    string
	children    []*route
	raw         raw
	middlewares []Middleware
}

func Route(items ...RouteItem) *route {
	r := &route{path: "/", middlewares: []Middleware{gziphandler.GzipHandler}}

	var resolveitem func(RouteItem)

	resolveitem = func(item RouteItem) {
		switch item := item.(type) {
		case *route:
			item.path = normalpath(path.Join(normalpath(r.path), normalpath(item.path)))
			if r.layout != nil {
				original := item.layout
				item.layout = func(lc *LayoutContext) h.I {
					if original != nil {
						lc.outlet = original(lc)
					}
					return r.layout(lc)
				}
			}
			r.children = append(r.children, item)
		case *public:
			r.public = item
		case *Fragment:
			for _, item := range item.items {
				resolveitem(item)
			}
		case Page:
			r.page = item
		case Layout:
			r.layout = item
		case Path:
			r.path = normalpath(string(item))
		case Redirect:
			r.redirect = string(item)
		case Middleware:
			r.middlewares = append(r.middlewares, item)
		default:
			panic(fmt.Sprintf("unknown item type %T", item))
		}
	}

	for _, item := range items {
		resolveitem(item)
	}

	return r
}

func Routes(items ...RouteItem) *route {
	r := Route(items...)
	r.root = true
	return r
}

type Path string

type Redirect string

type public struct {
	dir  fs.FS
	root string
}

func Public(dir fs.FS, root string) *public {
	return &public{dir: dir, root: root}
}

type Middleware func(http.Handler) http.Handler

type Fragment struct {
	items []RouteItem
}

func Frag(items ...RouteItem) *Fragment {
	return &Fragment{items: items}
}

func Map[T any](items []T, fn func(T, int) RouteItem) RouteItem {
	result := []RouteItem{}
	for i, item := range items {
		result = append(result, fn(item, i))
	}
	return Frag(result...)
}

func (*route) item()     {}
func (*public) item()    {}
func (*Fragment) item()  {}
func (Page) item()       {}
func (Layout) item()     {}
func (Path) item()       {}
func (Redirect) item()   {}
func (Middleware) item() {}

func default404(_ *PageContext) h.I {
	return h.Main(
		h.Class("w-full h-dvh flex justify-center items-center"),

		h.Div(
			h.Text("404 | Not Found"),
		),
	)
}

func (rt route) register(mux *http.ServeMux, head *c.AppHead) {
	if rt.root {
		mux.Handle("GET /ui/", gziphandler.GzipHandler(http.StripPrefix("/ui/", http.FileServerFS(head.FS()))))
	}

	if rt.public != nil {
		fs, err := fs.Sub(rt.public.dir, rt.public.root)
		if err != nil {
			panic(err)
		}
		mux.Handle("GET /public/", gziphandler.GzipHandler(http.StripPrefix("/public/", http.FileServerFS(fs))))
	}

	for _, child := range rt.children {
		child.middlewares = append(child.middlewares, rt.middlewares...)
		child.register(mux, head)
	}

	var handler http.Handler

	if rt.page != nil {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")

			ctx, cancel := context.WithCancel(r.Context())
			pctx := &PageContext{Context: ctx, cancel: cancel, req: r}

			var page Page
			if rt.path != r.URL.Path {
				w.WriteHeader(http.StatusNotFound)
				page = default404
			} else {
				w.WriteHeader(http.StatusOK)
				page = rt.page
			}

			var renderer h.I
			if rt.layout != nil {
				renderer = rt.layout(&LayoutContext{PageContext: pctx, outlet: page(pctx), head: head})
			} else {
				renderer = page(pctx)
			}

			renderer.Render(pctx, w)
		})
	} else if len(rt.redirect) != 0 {
		handler = http.RedirectHandler(rt.redirect, http.StatusFound)
	} else if len(rt.raw.bytes) != 0 {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", rt.raw.contenttyp)
			w.WriteHeader(200)
			w.Write(rt.raw.bytes)
		})
	}

	if handler != nil {
		for _, m := range rt.middlewares {
			handler = m(handler)
		}

		mux.Handle("GET "+rt.path, handler)
	}
}

func (rt route) Handler() http.Handler {
	mux := http.NewServeMux()
	head := c.CreateHead("/ui/")
	rt.register(mux, head)
	return mux
}

func Robots(robots []byte) RouteItem {
	return &route{
		path: "/robots.txt",
		raw: raw{
			bytes:      robots,
			contenttyp: "plain/text; charset=utf-8",
		},
	}
}

func Sitemap(sitemap []byte) RouteItem {
	return &route{
		path: "sitemap.xml",
		raw: raw{
			bytes:      sitemap,
			contenttyp: "application/xml",
		},
	}
}
