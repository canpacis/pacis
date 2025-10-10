package server

import (
	"context"
	"net/http"

	payload "github.com/canpacis/http-payload"
	"github.com/canpacis/pacis/html"
)

type chunk struct {
	id string
	fn func() *html.Element
}

type Context struct {
	context.Context

	chunks      []chunk
	speculation Speculation
}

func (c *Context) RegisterChunk(id string, fn func() *html.Element) {
	c.chunks = append(c.chunks, chunk{id: id, fn: fn})
}

func (c *Context) RegisterPrerender(url string) {
	if len(c.speculation.Prerender) == 0 {
		c.speculation.Prerender = append(c.speculation.Prerender, SpeculationRule{})
	}
	c.speculation.Prerender[0].URLs = append(c.speculation.Prerender[0].URLs, url)
}

func (c *Context) RegisterPrefetch(url string) {
	if len(c.speculation.Prefetch) == 0 {
		c.speculation.Prefetch = append(c.speculation.Prefetch, SpeculationRule{})
	}
	c.speculation.Prefetch[0].URLs = append(c.speculation.Prefetch[0].URLs, url)
}

type LayoutFn func(*App, context.Context, html.Node) html.Node

func HandlerOf[P any](app *App, fn func(context.Context, *P) html.Node, layout LayoutFn, middlewares ...func(http.Handler) http.Handler) http.Handler {
	middlewares = append(middlewares, app.middlewares...)

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Pattern == "/" || r.Pattern == "GET /") && r.URL.Path != "/" {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}

		p := new(P)
		ctx := &Context{
			Context: r.Context(),
		}

		scanner := payload.NewPipeScanner(
			payload.NewQueryScanner(r.URL.Query()),
			payload.NewHeaderScanner(&r.Header),
			payload.NewCookieScanner(r.Cookies()),
			payload.NewPathScanner(r),
		)

		var node html.Node
		var wrapper LayoutFn = layout
		if wrapper == nil {
			wrapper = func(app *App, ctx context.Context, n html.Node) html.Node { return n }
		}

		if err := scanner.Scan(p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			node = html.Error(err)
		} else {
			w.WriteHeader(http.StatusOK)
			node = fn(ctx, p)
		}
		app.speculation = ctx.speculation

		var flusher http.Flusher
		wrapper(app, ctx, node).Render(ctx, w)
		if len(ctx.chunks) != 0 {
			var ok bool
			flusher, ok = w.(http.Flusher)
			if !ok {
				// TODO: Maybe log?
				return
			}
		}

		if flusher != nil {
			flusher.Flush()

			for _, chunk := range ctx.chunks {
				element := chunk.fn()
				element.SetAttribute("slot", chunk.id)
				element.Render(ctx, w)
				flusher.Flush()
			}
		}
	})

	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

type handler struct{}

func BareHandlerOf(app *App, fn func(context.Context) html.Node, layout LayoutFn, middlewares ...func(http.Handler) http.Handler) http.Handler {
	return HandlerOf(app, func(ctx context.Context, p *handler) html.Node {
		return fn(ctx)
	}, layout, middlewares...)
}
