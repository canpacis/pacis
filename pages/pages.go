package pages

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/canpacis/pacis/pages/internal"
	"github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/scanner"
)

type apphead struct {
	content html.I
	meta    *Metadata
}

type Context struct {
	context.Context
	w http.ResponseWriter
	r *http.Request

	cookies []*http.Cookie
	pattern string

	head   apphead
	body   html.I
	outlet html.I
}

func (ctx *Context) Clone(parent context.Context) *Context {
	nctx := NewContext(ctx.w, ctx.r.Clone(parent))
	nctx.head = ctx.head
	nctx.body = ctx.body
	nctx.outlet = ctx.outlet
	nctx.cookies = ctx.cookies
	nctx.pattern = ctx.pattern
	ctx.Context = parent
	return nctx
}

func (ctx *Context) Request() *http.Request {
	return ctx.r
}

func (ctx *Context) Scan(v any) error {
	req := ctx.r.Clone(context.Background())
	query := req.URL.Query()

	pipe := scanner.NewPipe(
		scanner.NewCookie(ctx.cookies),
		scanner.NewQuery(&query),
		scanner.NewPath(req),
		scanner.NewHeader(&req.Header),
		internal.NewContextScanner(ctx),
	)
	return pipe.Scan(v)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w: w, r: r, Context: r.Context(), cookies: r.Cookies()}
}

type Page interface {
	Page(*Context) html.I
}

type FnPage struct {
	fn func(*Context) html.I
}

func (p FnPage) Page(ctx *Context) html.I {
	return p.fn(ctx)
}

func NewFnPage(fn func(*Context) html.I) *FnPage {
	return &FnPage{fn: fn}
}

type Layout interface {
	Layout(*Context) html.I
}

var EmptyLayout = NewFnLayout(func(ctx *Context) html.I {
	return ctx.outlet
})

type FnLayout struct {
	fn func(*Context) html.I
}

func (p FnLayout) Layout(ctx *Context) html.I {
	return p.fn(ctx)
}

func NewFnLayout(fn func(*Context) html.I) *FnLayout {
	return &FnLayout{fn: fn}
}

type Action interface {
	Action(*Context) html.I
}

type Route interface {
	http.Handler
	Path() string
}

func getdocrenderer(ctx *Context, pg Page, layout Layout, home bool) html.Renderer {
	var renderer html.I
	var page Page

	if home && ctx.r.URL.Path != "/" {
		ctx.w.WriteHeader(http.StatusNotFound)
		page = NotFoundPage
	} else {
		page = pg
	}
	_, ok := page.(*FnPage)
	if !ok {
		// If the page is not just a function page, create a new instance of the page type
		rv := reflect.ValueOf(page)
		rv = reflect.Indirect(rv)
		page = reflect.New(rv.Type()).Interface().(Page)
	}

	if err := ctx.Scan(page); err != nil {
		errpage := newerrpage()
		errpage.SetError(err)
		errpage.SetStatus(http.StatusInternalServerError)
		page = errpage
	}

	if layout != nil {
		_, ok := layout.(*FnLayout)
		if !ok {
			// If the layout is not just a function layout, create a new instance of the layout type
			rv := reflect.ValueOf(layout)
			rv = reflect.Indirect(rv)
			layout = reflect.New(rv.Type()).Interface().(Layout)
		}
		if err := ctx.Scan(layout); err != nil {
			errpage := newerrpage()
			errpage.SetError(err)
			errpage.SetStatus(http.StatusInternalServerError)
			page = errpage
		}

		ctx.outlet = page.Page(ctx)
		renderer = layout.Layout(ctx)
	} else {
		renderer = page.Page(ctx)
	}

	return renderer
}

type pageroute struct {
	path        string
	page        Page
	layout      Layout
	head        html.I
	body        html.I
	middlewares []func(http.Handler) http.Handler
}

func (r pageroute) Path() string {
	return r.path
}

func (pr *pageroute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		ctx.head = apphead{content: pr.head}
		ctx.body = pr.body
		ctx.pattern = strings.TrimPrefix(pr.path, "GET ")

		renderer := getdocrenderer(ctx, pr.page, pr.layout, pr.path == "GET /")
		sw := internal.NewStreamWriter(renderer, w)
		internal.Render(ctx, sw)
	})
	for _, middleware := range pr.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
}

func NewPageRoute(
	path string,
	page Page,
	layout Layout,
	head, body html.I,
	middlewares ...func(http.Handler) http.Handler,
) *pageroute {
	return &pageroute{
		path:   path,
		page:   page,
		layout: layout,
		head:   head, body: body,
		middlewares: middlewares,
	}
}

type actionroute struct {
	path        string
	action      Action
	middlewares []func(http.Handler) http.Handler
}

func (r actionroute) Path() string {
	return r.path
}

func (ar *actionroute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		renderer := ar.action.Action(ctx)
		sw := internal.NewStreamWriter(renderer, w)
		internal.Render(ctx, sw)
	})
	for _, middleware := range ar.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
}

func NewActionRoute(path string, action Action, middlewares ...func(http.Handler) http.Handler) *actionroute {
	return &actionroute{path: path, action: action, middlewares: middlewares}
}

type redirectroute struct {
	path        string
	to          string
	code        int
	middlewares []func(http.Handler) http.Handler
}

func NewRedirectRoute(path, to string, code int, middlewares ...func(http.Handler) http.Handler) *redirectroute {
	return &redirectroute{path: path, to: to, code: code, middlewares: middlewares}
}

func (rr redirectroute) Path() string {
	return rr.path
}

func (rr *redirectroute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.Redirect(w, r, rr.to, rr.code)
	})
	for _, middleware := range rr.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
}

type rawroute struct {
	path        string
	contenttyp  string
	content     []byte
	middlewares []func(http.Handler) http.Handler
}

func NewRawRoute(path, typ string, content []byte, middlewares ...func(http.Handler) http.Handler) *rawroute {
	return &rawroute{path: path, contenttyp: typ, content: content, middlewares: middlewares}
}

func (rr rawroute) Path() string {
	return rr.path
}

func (rr *rawroute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func WrapLayout(layout Layout, rest ...Layout) Layout {
	switch len(rest) {
	case 0:
		return layout
	case 1:
		return NewFnLayout(func(ctx *Context) html.I {
			_, ok := layout.(*FnLayout)
			if !ok {
				ctx.Scan(layout)
			}
			_, ok = rest[0].(*FnLayout)
			if !ok {
				ctx.Scan(rest[0])
			}
			ctx.outlet = layout.Layout(ctx)
			return rest[0].Layout(ctx)
		})
	default:
		first := NewFnLayout(func(ctx *Context) html.I {
			_, ok := layout.(*FnLayout)
			if !ok {
				ctx.Scan(layout)
			}
			_, ok = rest[0].(*FnLayout)
			if !ok {
				ctx.Scan(rest[0])
			}
			ctx.outlet = layout.Layout(ctx)
			return rest[0].Layout(ctx)
		})
		return WrapLayout(first, rest[1:]...)
	}
}

var NotFoundPage Page = &FnPage{func(*Context) html.I {
	return html.P(html.Text("Not Found"))
}}

func SetNotFoundPage(page Page) {
	NotFoundPage = page
}

type ErrorPage interface {
	error
	Page
	SetError(error)
	Status() int
	SetStatus(int)
}

type DefaultErrorPage struct {
	err    error
	status int
}

func (dep *DefaultErrorPage) Error() string {
	return dep.err.Error()
}

func (dep *DefaultErrorPage) SetError(err error) {
	dep.err = err
}

func (dep *DefaultErrorPage) Page(ctx *Context) html.I {
	return html.P(html.Text(dep.Error()))
}

func (dep *DefaultErrorPage) Status() int {
	return dep.status
}

func (dep *DefaultErrorPage) SetStatus(status int) {
	dep.status = status
}

var deferrpage ErrorPage = &DefaultErrorPage{}

func newerrpage() ErrorPage {
	return reflect.New(reflect.Indirect(reflect.ValueOf(deferrpage)).Type()).Interface().(ErrorPage)
}

func SetErrorPage(page ErrorPage) {
	deferrpage = page
}

func Serve(addr string, router http.Handler, logger *slog.Logger) {
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if logger == nil {
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
	slog.SetDefault(logger)

	go func() {
		logger.Debug("Server is starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}
	logger.Debug("Server stopped")
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
