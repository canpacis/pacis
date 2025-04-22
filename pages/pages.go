package pages

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/scanner"
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

func SafeGet[T any](ctx context.Context, key string) (T, bool) {
	value := ctx.Value(ctxkey(fmt.Sprintf("%s:%s", "app", key)))
	cast, ok := value.(T)
	return cast, ok
}

type PageContext struct {
	context.Context
	w       http.ResponseWriter
	r       *http.Request
	hookch  chan string
	chsize  atomic.Int32
	elemch  chan h.Element
	ready   atomic.Bool
	timings []*ServerTiming
	logger  *slog.Logger
	title   string
}

// func (ctx *PageContext) Message() {

// }

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

func (ctx *PageContext) Error(status int) h.I {
	ctx.w.Header().Set("Content-Type", "text/html")
	ctx.w.WriteHeader(status)

	return ErrorPage(ctx)
}

func (ctx *PageContext) Logger() *slog.Logger {
	return ctx.logger
}

func (ctx *PageContext) Scan(v any) error {
	query := ctx.r.URL.Query()

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	jar.SetCookies(ctx.r.URL, ctx.r.Cookies())

	pipe := scanner.NewPipe(
		scanner.NewQuery(&query),
		scanner.NewHeader(&ctx.r.Header),
		scanner.NewCookie(jar, ctx.r.URL),
	)
	return pipe.Scan(v)
}

func (ctx *PageContext) Set(key string, value any) {
	c := context.WithValue(ctx.Context, ctxkey(fmt.Sprintf("%s:%s", "app", key)), value)
	ctx.Context = c
}

func (ctx *PageContext) GetCookie(name string) (*http.Cookie, error) {
	return ctx.r.Cookie(name)
}

func (ctx *PageContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.w, cookie)
}

// TODO: Replace this API with Page Meta Data API
func (ctx *PageContext) SetTitle(title string) {
	ctx.title = title
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

type ActionContext struct {
	*PageContext
}

func (ctx *ActionContext) Scan(v any) error {
	if err := ctx.r.ParseForm(); err != nil {
		return err
	}

	query := ctx.r.URL.Query()

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	jar.SetCookies(ctx.r.URL, ctx.r.Cookies())

	pipe := scanner.NewPipe(
		scanner.NewQuery(&query),
		scanner.NewForm(&ctx.r.PostForm),
		scanner.NewHeader(&ctx.r.Header),
		scanner.NewCookie(jar, ctx.r.URL),
	)

	return pipe.Scan(v)
}

type Action func(*ActionContext) h.I

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

var ErrorPage Page = func(pc *PageContext) h.I {
	return h.P(h.Text("Error"))
}

func SetErrorPage(page Page) {
	ErrorPage = page
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

type Streamer struct {
	layout     Layout
	page       Page
	head, body h.I
	basecheck  bool

	limit int
	buf   *bytes.Buffer
	w     http.ResponseWriter
	http.Flusher
}

func (s *Streamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	s.w = w
	s.Flusher = w.(http.Flusher)
	s.buf.Reset()

	ctx := &PageContext{w: w, r: r, logger: slog.Default(), Context: r.Context()}
	var renderer h.I
	var page Page

	if s.basecheck {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			page = NotFoundPage
		} else {
			w.WriteHeader(http.StatusOK)
			page = s.page
		}
	} else {
		page = s.page
	}

	if s.layout != nil {
		outlet := page(ctx)

		// TODO: Maybe do the page metadata api with hooks?
		if len(ctx.title) != 0 {
			el, ok := s.head.(h.Element)
			if ok {
				el.AddNode(h.Title(h.Text(ctx.title)))
			} else {
				frag, ok := s.head.(*h.Fragment)
				if ok {
					s.head = h.Frag(frag, h.Title(h.Text(ctx.title)))
				} else {
					ctx.logger.Error("failed to add title to document head")
				}
			}
		}

		renderer = s.layout(&LayoutContext{PageContext: ctx, head: s.head, body: s.body, outlet: outlet})
	} else {
		renderer = page(ctx)
	}

	renderer.Render(ctx, s)

	// TODO: Maybe carry this to the hooks api
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
			return
		case el := <-ctx.elemch:
			el.Render(ctx, w)
			s.Flush()
		}
	}
}

func (s *Streamer) Write(p []byte) (int, error) {
	if s.buf.Len() < s.limit {
		return s.buf.Write(p)
	}

	n, err := s.buf.Write(p)
	if err != nil {
		return n, err
	}
	m, err := io.Copy(s.w, s.buf)
	s.Flush()
	s.buf.Reset()
	return int(m), err
}

func NewStreamer(layout Layout, page Page, head, body h.I, basecheck bool) *Streamer {
	return &Streamer{
		layout: layout,
		page:   page,
		head:   head, body: body,
		basecheck: basecheck,
		buf:       new(bytes.Buffer),
		limit:     500,
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
	var handler http.Handler = NewStreamer(hr.layout, hr.page, hr.head, hr.body, true)
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
	var handler http.Handler = NewStreamer(pr.layout, pr.page, pr.head, pr.body, false)
	for _, middleware := range pr.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
}

type RedirectRoute struct {
	path        string
	to          string
	code        int
	middlewares []func(http.Handler) http.Handler
}

func NewRedirectRoute(path, to string, code int, middlewares ...func(http.Handler) http.Handler) *RedirectRoute {
	return &RedirectRoute{path: path, to: to, code: code, middlewares: middlewares}
}

func (rr RedirectRoute) Path() string {
	return rr.path
}

func (rr *RedirectRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, rr.to, rr.code)
	})
	for _, middleware := range rr.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
}

type ActionRoute struct {
	path        string
	action      Action
	middlewares []func(http.Handler) http.Handler
}

func NewActionRoute(path string, action Action, middlewares ...func(http.Handler) http.Handler) *ActionRoute {
	return &ActionRoute{path: path, action: action, middlewares: middlewares}
}

func (ar ActionRoute) Path() string {
	return ar.path
}

func (ar *ActionRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = NewStreamer(EmptyLayout, func(ctx *PageContext) h.I {
		return ar.action(&ActionContext{PageContext: ctx})
	}, h.Frag(), h.Frag(), false)
	for _, middleware := range ar.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r)
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

func EmptyLayout(ctx *LayoutContext) h.I {
	return ctx.Outlet()
}
