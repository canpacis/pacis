package server

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/server/middleware"
)

// LayoutFn defines a function type that takes a context and an html.Node as input,
// and returns a modified html.Node. It is typically used to apply layout transformations
// or wrappers to HTML nodes within a given context.
type LayoutFn func(*App, html.Node) html.Node

type async struct {
	id   string
	comp html.Component
}

type redirect struct {
	status int
	to     string
}

type serverctx struct {
	context.Context

	async    []async
	redirect *redirect
	notfound bool
	req      *http.Request
}

/*
HandlerOf creates an HTTP handler that processes requests using the provided function `fn`,
which takes a context and a pointer to a parameter struct of type P. The handler applies the
specified layout function `layout` to the resulting HTML node, and supports an optional list
of middleware functions. The handler automatically scans and populates the parameter struct
from the request's query parameters, headers, cookies, and path variables. If scanning fails,
it responds with a 400 Bad Request and an error node. Otherwise, it renders the node returned
by `fn`, applies the layout, and supports chunked rendering if needed. All provided and
application-level middlewares are applied in order.

Type Parameters:
  - P: The type of the parameter struct to be populated from the request.

Parameters:
  - app: The application context containing shared resources and middlewares.
  - fn: A function that generates an HTML node given a context and a pointer to P.
  - layout: A layout function to wrap the generated HTML node. If nil, a default passthrough is used.
  - middlewares: Optional HTTP middleware functions to wrap the handler.

Returns:
  - http.Handler: The composed HTTP handler ready to be registered with a router or server.
*/
func HandlerOf(app *App, fn func() html.Node, layout LayoutFn, middlewares ...middleware.Middleware) http.Handler {
	var wrapper LayoutFn = layout
	if wrapper == nil {
		wrapper = func(app *App, n html.Node) html.Node { return n }
	}
	node := wrapper(app, fn())

	renderer := NewStaticRenderer()
	if err := renderer.Build(node); err != nil {
		log.Fatalf("Failed to statically render page: %s", err.Error())
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Pattern == "/" || r.Pattern == "GET /") && r.URL.Path != "/" {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}

		ctx := &serverctx{Context: r.Context(), req: r}

		if ctx.redirect != nil {
			w.WriteHeader(ctx.redirect.status)
			http.Redirect(w, r, ctx.redirect.to, ctx.redirect.status)
		} else if ctx.notfound {
			w.WriteHeader(http.StatusNotFound)
			http.NotFoundHandler().ServeHTTP(w, r)
		} else {
			bw := bufio.NewWriter(w)

			if err := renderer.Render(ctx, bw); err != nil {
				return
			}

			bw.Flush()

			flusher, ok := w.(http.Flusher)
			if !ok {
				log.Fatal("http writer does not support chunked encoding")
				return
			}
			if len(ctx.async) == 0 {
				return
			}
			flusher.Flush()

			renderers := make(chan *StaticRenderer)
			wg := sync.WaitGroup{}

			for _, chunk := range ctx.async {
				wg.Add(1)
				go func(chunk async) {
					defer wg.Done()

					renderer := NewStaticRenderer()
					var node html.Node
					node = chunk.comp(ctx)
					elem, ok := node.(*html.Element)
					if !ok {
						node = html.Template(html.SlotAttr(chunk.id), node)
					} else {
						elem.SetAttribute("slot", chunk.id)
					}

					if err := renderer.Build(node); err != nil {
						log.Fatalf("Failed to statically render page: %s", err.Error())
						return
					}

					renderers <- renderer
				}(chunk)
			}

			go func() {
				wg.Wait()
				close(renderers)
			}()

			for renderer := range renderers {
				renderer.Render(ctx, w)
				flusher.Flush()
				renderer.Release()
			}
		}
	})

	for i := len(app.middlewares) - 1; i >= 0; i-- {
		handler = app.middlewares[i].Apply(handler)
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Apply(handler)
	}

	return handler
}
