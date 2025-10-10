package server

import (
	"context"
	"net/http"

	payload "github.com/canpacis/http-payload"
	"github.com/canpacis/pacis/html"
)

/*
# Speculation Rules API

A string providing a hint to the browser as to how eagerly it should prefetch/prerender link targets in order to balance performance advantages against resource overheads.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#eagerness
*/
type SpeculationEagerness string

const (
	/*
		# Speculation Rules API

		The author thinks the link is very likely to be followed, and/or the document may take significant time to fetch. Prefetch/prerender should start as soon as possible, subject only to considerations such as user preferences and resource limits.

		https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#immediate
	*/
	ImmediateEagerness = SpeculationEagerness("immediate")
	/*
		# Speculation Rules API

		The author wants to prefetch/prerender a large number of navigations, as early as possible. Prefetch/prerender should start on any slight suggestion that a link may be followed. For example, the user could move their mouse cursor towards the link, hover/focus it for a moment, or pause scrolling with the link in a prominent place.

		https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#eager
	*/
	EagerEagerness = SpeculationEagerness("eager")
	/*
		# Speculation Rules API

		The author is looking for a balance between eager and conservative. Prefetch/prerender should start when there is a reasonable suggestion that the user will follow a link in the near future. For example, the user could scroll a link into the viewport and hover/focus it for some time.

		https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#moderate
	*/
	ModerateEagerness = SpeculationEagerness("moderate")
	/*
		# Speculation Rules API

		The author wishes to get some benefit from speculative loading with a fairly small tradeoff of resources. Prefetch/prerender should start only when the user is starting to click on the link, for example on mousedown or pointerdown.

		https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#conservative
	*/
	ConservativeEagerness = SpeculationEagerness("conservative")
)

/*
# Speculation Rules API

The speculationrules value of the type attribute of the <script> element indicates that the body of the element contains speculation rules.

Speculation rules take the form of a JSON structure that determine what resources should be prefetched or prerendered by the browser. This is part of the Speculation Rules API.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules
*/
type SpeculationRule struct {
	URLs      []string             `json:"urls,omitempty"`
	Eagerness SpeculationEagerness `json:"eagerness"`
}

/*
# Speculation Rules API

The JSON structure contains one or more fields at the top level, each one representing an action to define speculation rules for. At present the supported actions are:

"prefetch" Optional Experimental
Rules for potential future navigations that should have their associated document response body downloaded, leading to significant performance improvements when those documents are navigated to. Note that none of the subresources referenced by the page are downloaded.

"prerender" Optional Experimental
Rules for potential future navigations that should have their associated documents fully downloaded, rendered, and loaded into an invisible tab. This includes loading all subresources, running all JavaScript, and even loading subresources and performing data fetches started by JavaScript. When those documents are navigated to, navigations will be instant, leading to major performance improvements.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules
*/
type Speculation struct {
	Prerender []SpeculationRule `json:"prerender,omitempty"`
	Prefetch  []SpeculationRule `json:"prefetch,omitempty"`
}

type chunk struct {
	id string
	fn func() *html.Element
}

// Context embeds context.Context and provides additional fields for application-specific data.
// It holds a reference to the App, a slice of chunk objects, and Speculation data for request handling.
type Context struct {
	context.Context

	app    *App
	chunks []chunk
	specs  Speculation
}

// RegisterChunk registers a new chunk with the given identifier and a function that returns an *html.Element.
// The chunk is appended to the context's list of chunks for later processing or rendering.
//
// Parameters:
//   - id: A unique string identifier for the chunk.
//   - fn: A function that returns a pointer to an html.Element representing the chunk's content.
func (c *Context) RegisterChunk(id string, fn func() *html.Element) {
	c.chunks = append(c.chunks, chunk{id: id, fn: fn})
}

// RegisterSpeculation adds a SpeculationRule to either the Prerender or Prefetch list
// in the Context's specs, depending on the value of the render parameter.
// If render is true, the rule is added to the Prerender list; otherwise, it is added to the Prefetch list.
//
// Parameters:
//   - rule:   The SpeculationRule to be registered.
//   - render: A boolean indicating whether to add the rule to Prerender (true) or Prefetch (false).
func (c *Context) RegisterSpeculation(rule SpeculationRule, render bool) {
	if render {
		c.specs.Prerender = append(c.specs.Prerender, rule)
	} else {
		c.specs.Prefetch = append(c.specs.Prefetch, rule)
	}
}

// LayoutFn defines a function type that takes a context and an html.Node as input,
// and returns a modified html.Node. It is typically used to apply layout transformations
// or wrappers to HTML nodes within a given context.
type LayoutFn func(context.Context, html.Node) html.Node

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
			app:     app,
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
			wrapper = func(ctx context.Context, n html.Node) html.Node { return n }
		}

		if err := scanner.Scan(p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			node = html.Error(err)
		} else {
			w.WriteHeader(http.StatusOK)
			node = fn(ctx, p)
		}

		var flusher http.Flusher
		wrapper(ctx, node).Render(ctx, w)
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

/*
BareHandlerOf creates an HTTP handler that processes requests using the provided function `fn`,
which takes only a context. The handler applies the specified layout function `layout` to the
resulting HTML node, and supports an optional list of middleware functions. The handler
automatically scans and populates the parameter struct from the request's query parameters,
headers, cookies, and path variables. If scanning fails, it responds with a 400 Bad Request
and an error node. Otherwise, it renders the node returned by `fn`, applies the layout, and
supports chunked rendering if needed. All provided and application-level middlewares are applied
in order.

Parameters:
  - app: The application context containing shared resources and middlewares.
  - fn: A function that generates an HTML node given a context.
  - layout: A layout function to wrap the generated HTML node. If nil, a default passthrough is used.
  - middlewares: Optional HTTP middleware functions to wrap the handler.

Returns:
  - http.Handler: The composed HTTP handler ready to be registered with a router or server.
*/
func BareHandlerOf(app *App, fn func(context.Context) html.Node, layout LayoutFn, middlewares ...func(http.Handler) http.Handler) http.Handler {
	return HandlerOf(app, func(ctx context.Context, p *handler) html.Node {
		return fn(ctx)
	}, layout, middlewares...)
}
