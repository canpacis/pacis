package pages

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
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
	w       http.ResponseWriter
	r       *http.Request
	chsize  atomic.Int32
	elemch  chan h.Element
	ready   atomic.Bool
	timings []*ServerTiming
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

func (ctx *PageContext) QueueElement() func(h.Element) {
	ctx.chsize.Add(1)
	return func(el h.Element) {
		ctx.elemch <- el
	}
}

func (ctx *PageContext) Ready() bool {
	return ctx.ready.Load()
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

func (ctx *PageContext) AddTiming(timing *ServerTiming) {
	ctx.timings = append(ctx.timings, timing)
}

func (ctx *PageContext) Set(key string, value any) {
	c := context.WithValue(ctx, ctxkey(fmt.Sprintf("%s:%s", "app", key)), value)
	ctx.r = ctx.r.Clone(c)
}

type Page func(*PageContext) h.I

type LayoutContext struct {
	*PageContext
	head   h.I
	body   h.I
	outlet h.I
}

func (ctx LayoutContext) Head() h.I {
	return ctx.head
}

func (ctx LayoutContext) Body() h.I {
	return ctx.body
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

func render(item h.I, w http.ResponseWriter, ctx *PageContext, timing *Timing) {
	flusher := w.(http.Flusher)
	buf := new(bytes.Buffer)

	item.Render(ctx, buf)
	timing.Done(ctx)

	timings := []string{}
	for _, timing := range ctx.timings {
		timings = append(timings, timing.String())
	}
	w.Header().Set("Server-Timing", strings.Join(timings, ","))

	io.Copy(w, buf)
	flusher.Flush()

	size := int(ctx.chsize.Load())
	if size == 0 {
		return
	}

	ctx.elemch = make(chan h.Element, size)
	ctx.ready.Store(true)

	for range size {
		select {
		case <-ctx.Done():
			// client disconnected
		case el := <-ctx.elemch:
			el.Render(ctx, w)
			flusher.Flush()
		}
	}
}

type Route interface {
	http.Handler
	Path() string
}

type HomeRoute struct {
	page        Page
	layout      Layout
	head        h.I
	body        h.I
	middlewares []func(http.Handler) http.Handler
}

func (HomeRoute) Path() string {
	return "/"
}

func (hr *HomeRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		rt := NewTiming("render", "Document rendered")

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
			pt := NewTiming("page", "Page prepared")
			outlet := page(ctx)
			pt.Done(ctx)
			lt := NewTiming("layout", "Layout prepared")
			renderer = hr.layout(&LayoutContext{PageContext: ctx, head: hr.head, body: hr.body, outlet: outlet})
			lt.Done(ctx)
		} else {
			pt := NewTiming("page", "Page prepared")
			renderer = page(ctx)
			pt.Done(ctx)
		}

		render(renderer, w, ctx, rt)
	})
	for _, middleware := range hr.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
}

func NewHomeRoute(page Page, layout Layout, head, body h.I, middlewares ...func(http.Handler) http.Handler) *HomeRoute {
	return &HomeRoute{page: page, layout: layout, head: head, body: body, middlewares: middlewares}
}

type PageRoute struct {
	path        string
	page        Page
	layout      Layout
	head        h.I
	body        h.I
	middlewares []func(http.Handler) http.Handler
}

func NewPageRoute(path string, page Page, layout Layout, head, body h.I, middlewares ...func(http.Handler) http.Handler) *PageRoute {
	return &PageRoute{path: path, page: page, layout: layout, head: head, body: body, middlewares: middlewares}
}

func (pr PageRoute) Path() string {
	return pr.path
}

func (pr *PageRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		rt := NewTiming("render", "Document rendered")

		ctx := &PageContext{w: w, r: r}
		var renderer h.I
		if pr.layout != nil {
			pt := NewTiming("page", "Page prepared")
			outlet := pr.page(ctx)
			pt.Done(ctx)
			lt := NewTiming("layout", "Layout prepared")
			renderer = pr.layout(&LayoutContext{PageContext: ctx, head: pr.head, body: pr.body, outlet: outlet})
			lt.Done(ctx)
		} else {
			pt := NewTiming("page", "Page prepared")
			renderer = pr.page(ctx)
			pt.Done(ctx)
		}

		render(renderer, w, ctx, rt)
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

func Serve(addr string, router http.Handler) {
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
