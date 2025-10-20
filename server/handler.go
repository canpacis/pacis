package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/canpacis/pacis/html"
	intserver "github.com/canpacis/pacis/internal"
	"github.com/canpacis/pacis/server/metadata"
	"github.com/canpacis/pacis/server/middleware"
)

type Page interface {
	Page() html.Node
}

type PageFunc func() html.Node

func (p PageFunc) Page() html.Node {
	return p()
}

func (PageFunc) Metadata() *metadata.Metadata {
	return &metadata.Metadata{}
}

// Layout defines a function type that takes a context and an html.Node as input,
// and returns a modified html.Node. It is typically used to apply layout transformations
// or wrappers to HTML nodes within a given context.
type Layout func(*Server, html.Node, html.Node) html.Node

func DefaultLayout(server *Server, head html.Node, children html.Node) html.Node {
	return html.Fragment(
		html.Doctype,
		html.Head(
			html.Meta(html.Charset("UTF-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1.0")),
			head,
		),
		html.Body(
			children,
		),
	)
}

func head(page Page, dev bool) html.Node {
	staticmeta, ok := page.(interface{ Metadata() *metadata.Metadata })
	if ok {
		if dev {
			return html.Fragment(
				staticmeta.Metadata().Node(),
				html.Script(html.Type("module"), html.Src("/@vite/client")),
			)
		}
		return staticmeta.Metadata().Node()
	} else {
		dynamicmeta, ok := page.(interface {
			Metadata(context.Context) *metadata.Metadata
		})
		if !ok {
			log.Fatalf("Invalid page type %T, type must have a `Metadata() *metadata.Metadata` method to implement the Page interface.", page)
		}
		if dev {
			return html.Fragment(
				html.Component(func(ctx context.Context) html.Node {
					return dynamicmeta.Metadata(ctx).Node()
				}),
				html.Script(html.Type("module"), html.Src("/@vite/client")),
			)
		}
		return html.Component(func(ctx context.Context) html.Node {
			return dynamicmeta.Metadata(ctx).Node()
		})
	}
}

var bufpool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

/*
PageHandler creates an HTTP handler that processes requests using the provided function `fn`,
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
func PageHandler(server *Server, page Page, layout Layout, middlewares ...middleware.Middleware) http.Handler {
	var handler = handler(server, page, layout, false)

	for i := len(server.middlewares) - 1; i >= 0; i-- {
		handler = server.middlewares[i].Apply(handler)
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Apply(handler)
	}

	return handler
}

func handler(server *Server, page Page, layout Layout, internal bool) http.Handler {
	defer func() {
		if data := recover(); data != nil {
			server.options.Logger.Error("HTTP handler paniced on partial pre-render", "error", data)
		}
	}()

	wrapper := layout
	if wrapper == nil {
		wrapper = func(s *Server, h, c html.Node) html.Node { return c }
	}
	node := wrapper(server, head(page, server.options.Env == Dev), page.Page())

	renderer := NewStaticRenderer()
	if err := renderer.Build(node); err != nil {
		log.Fatalf("Failed to statically render page: %s", err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := intserver.NewContext(w, r)

		buf := bufpool.New().(*bytes.Buffer)
		defer bufpool.Put(buf)

		if err := renderer.Render(ctx, buf); err != nil {
			return
		}

		if ctx.NotFoundMark {
			server.notfound.ServeHTTP(w, r)
			return
		}
		if ctx.RedirectMark != nil {
			http.Redirect(w, r, ctx.RedirectMark.To, ctx.RedirectMark.Status)
			return
		}

		if internal {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		io.Copy(w, buf)

		flusher, ok := w.(http.Flusher)
		if !ok {
			server.options.Logger.Warn("Http writer does not support chunked encoding")
			return
		}
		if len(ctx.AsyncChunks) == 0 {
			return
		}
		flusher.Flush()

		renderers := make(chan *StaticRenderer)
		wg := sync.WaitGroup{}

		for _, chunk := range ctx.AsyncChunks {
			wg.Add(1)
			go func(chunk intserver.AsyncChunk) {
				defer wg.Done()

				renderer := NewStaticRenderer()
				var node html.Node
				node = chunk.Component(ctx)
				elem, ok := node.(*html.Element)
				if !ok {
					node = html.Template(html.SlotAttr(chunk.ID), node)
				} else {
					elem.SetAttribute("slot", chunk.ID)
				}

				if err := renderer.Build(node); err != nil {
					server.options.Logger.Error("Failed to statically render page", "error", err)
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
		}
	})
}

type StaticRenderer struct {
	chunks []any
}

func (r *StaticRenderer) Build(node html.Node) error {
	buf := new(bytes.Buffer)
	cw := html.NewChunkWriter()
	node.Render(cw)
	defer node.Release()

	for _, chunk := range cw.Chunks() {
		switch chunk := chunk.(type) {
		case html.StaticChunk:
			if _, err := buf.Write(chunk); err != nil {
				return err
			}
		case html.DynamicChunk:
			r.chunks = append(r.chunks, buf.Bytes())
			bufpool.Put(buf)
			buf = bufpool.New().(*bytes.Buffer)
			r.chunks = append(r.chunks, chunk)
		default:
			return fmt.Errorf("invalid chunk type %T", chunk)
		}
	}

	r.chunks = append(r.chunks, buf.Bytes())
	return nil
}

func (r *StaticRenderer) Render(ctx context.Context, w io.Writer) error {
	for _, chunk := range r.chunks {
		switch chunk := chunk.(type) {
		case []byte:
			if _, err := w.Write(chunk); err != nil {
				return err
			}
		case html.DynamicChunk:
			if err := chunk(ctx, w); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid chunk type %t", chunk)
		}
	}
	return nil
}

func (r *StaticRenderer) Clear() {
	r.chunks = []any{}
}

func NewStaticRenderer() *StaticRenderer {
	return &StaticRenderer{
		chunks: []any{},
	}
}

type ActionFunc func(http.ResponseWriter, *http.Request) error

func ActionsHandler(server *Server, actions map[string]ActionFunc, middlewares ...middleware.Middleware) http.Handler {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("__action")
		action, ok := actions[name]
		if !ok {
			w.Write([]byte("unknown action"))
			return
		}
		err := action(w, r)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	})

	for i := len(server.middlewares) - 1; i >= 0; i-- {
		handler = server.middlewares[i].Apply(handler)
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Apply(handler)
	}
	return handler
}
