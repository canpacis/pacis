package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/canpacis/pacis/html"
	intserver "github.com/canpacis/pacis/internal/server"
	"github.com/canpacis/pacis/server/middleware"
)

// LayoutFn defines a function type that takes a context and an html.Node as input,
// and returns a modified html.Node. It is typically used to apply layout transformations
// or wrappers to HTML nodes within a given context.
type LayoutFn func(*Server, html.Node) html.Node

type PageFn func(*Server) html.Node

var bufpool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
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
func HandlerOf(server *Server, page PageFn, layout LayoutFn, middlewares ...middleware.Middleware) http.Handler {
	var wrapper LayoutFn = layout
	if wrapper == nil {
		wrapper = func(app *Server, n html.Node) html.Node { return n }
	}
	node := wrapper(server, page(server))

	renderer := NewStaticRenderer()
	if err := renderer.Build(node); err != nil {
		log.Fatalf("Failed to statically render page: %s", err.Error())
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Pattern == "/" || r.Pattern == "GET /") && r.URL.Path != "/" {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}

		ctx := intserver.NewContext(w, r)

		buf := bufpool.New().(*bytes.Buffer)
		defer bufpool.Put(buf)

		if err := renderer.Render(ctx, buf); err != nil {
			return
		}

		if ctx.NotFoundMark {
			w.WriteHeader(http.StatusNotFound)
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}
		if ctx.RedirectMark != nil {
			http.Redirect(w, r, ctx.RedirectMark.To, ctx.RedirectMark.Status)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.Copy(w, buf)

		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Fatal("http writer does not support chunked encoding")
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

type StaticRenderer struct {
	chunks []any
}

func (r *StaticRenderer) Build(node html.Node) error {
	buf := bufpool.New().(*bytes.Buffer)
	defer bufpool.Put(buf)

	for chunk := range node.Chunks() {
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
	bw := bufio.NewWriter(w)

	for _, chunk := range r.chunks {
		switch chunk := chunk.(type) {
		case []byte:
			if _, err := bw.Write(chunk); err != nil {
				return err
			}
		case html.DynamicChunk:
			if err := chunk(ctx, bw); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid chunk type %t", chunk)
		}
	}
	return bw.Flush()
}

func (r *StaticRenderer) Clear() {
	r.chunks = []any{}
}

func NewStaticRenderer() *StaticRenderer {
	return &StaticRenderer{
		chunks: []any{},
	}
}
