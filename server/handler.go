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
	intserver "github.com/canpacis/pacis/internal/server"
	"github.com/canpacis/pacis/server/metadata"
	"github.com/canpacis/pacis/server/middleware"
)

type layout struct {
	layout func(*Server, html.Node, html.Node) html.Node
}

func layoutof(l any) *layout {
	if l == nil {
		return &layout{
			layout: func(s *Server, h, c html.Node) html.Node {
				return c
			},
		}
	}
	shc, ok := l.(func(*Server, html.Node, html.Node) html.Node)
	if ok {
		return &layout{layout: shc}
	}
	hc, ok := l.(func(html.Node, html.Node) html.Node)
	if ok {
		return &layout{
			layout: func(s *Server, h, c html.Node) html.Node {
				return hc(h, c)
			},
		}
	}
	c, ok := l.(func(html.Node) html.Node)
	if ok {
		return &layout{
			layout: func(s *Server, h, ch html.Node) html.Node {
				return c(ch)
			},
		}
	}
	log.Fatal("Invalid layout type, you must pass a function that returns an html.Node.")
	return nil
}

// LayoutFn defines a function type that takes a context and an html.Node as input,
// and returns a modified html.Node. It is typically used to apply layout transformations
// or wrappers to HTML nodes within a given context.
type LayoutFn any

type page struct {
	staticmeta bool
	metadata   func(context.Context) *metadata.Metadata
	page       func(*Server) html.Node
}

func (p *page) Metadata(ctx context.Context) *metadata.Metadata {
	return p.metadata(ctx)
}

func (p *page) Page(s *Server) html.Node {
	return p.page(s)
}

func pageof(p any) *page {
	fn, ok := p.(func() html.Node)
	if ok {
		return &page{
			metadata: func(ctx context.Context) *metadata.Metadata {
				return nil
			},
			page: func(s *Server) html.Node {
				return fn()
			},
		}
	}

	serverfn, ok := p.(func(*Server) html.Node)
	if ok {
		return &page{
			metadata: func(ctx context.Context) *metadata.Metadata {
				return nil
			},
			page: func(s *Server) html.Node {
				return serverfn(s)
			},
		}
	}

	var staticmeta bool
	var metadatafn func(context.Context) *metadata.Metadata
	var pagefn func(*Server) html.Node

	metadataiface, ok := p.(interface{ Metadata() *metadata.Metadata })
	if ok {
		staticmeta = true
		metadatafn = func(ctx context.Context) *metadata.Metadata {
			return metadataiface.Metadata()
		}
	} else {
		ctxmetadataiface, ok := p.(interface {
			Metadata(context.Context) *metadata.Metadata
		})
		if ok {
			metadatafn = func(ctx context.Context) *metadata.Metadata {
				return ctxmetadataiface.Metadata(ctx)
			}
		} else {
			staticmeta = true
			metadatafn = func(ctx context.Context) *metadata.Metadata {
				return nil
			}
		}
	}

	pageiface, ok := p.(interface{ Page() html.Node })
	if ok {
		pagefn = func(s *Server) html.Node {
			return pageiface.Page()
		}
	} else {
		serverpageiface, ok := p.(interface{ Page(*Server) html.Node })
		if ok {
			pagefn = func(s *Server) html.Node {
				return serverpageiface.Page(s)
			}
		} else {
			log.Fatal("Invalid page type, you must either pass a function or an interface with a Page() method on it.")
		}
	}

	return &page{
		staticmeta: staticmeta,
		metadata:   metadatafn,
		page:       pagefn,
	}
}

type Page any

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
func HandlerOf(server *Server, page Page, layout LayoutFn, middlewares ...middleware.Middleware) http.Handler {
	defer func() {
		if data := recover(); data != nil {
			server.options.Logger.Error("HTTP handler paniced on partial pre-render", "error", data)
		}
	}()

	wrapper := layoutof(layout)
	p := pageof(page)
	var metanode html.Node
	if p.staticmeta {
		metanode = p.Metadata(context.Background()).Node()
	} else {
		metanode = html.Component(func(ctx context.Context) html.Node {
			return p.Metadata(ctx).Node()
		})
	}
	node := wrapper.layout(server, metanode, p.Page(server))

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
		ctx.Metadata = p.Metadata(ctx)

		buf := bufpool.New().(*bytes.Buffer)
		defer bufpool.Put(buf)

		if err := renderer.Render(ctx, buf); err != nil {
			return
		}

		if ctx.NotFoundMark {
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
