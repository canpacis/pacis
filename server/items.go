package server

import (
	"context"
	"io"
	"log"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/internal/util"
)

// Asset returns the URL or path for a given asset name based on the current application context.
// In development mode, it constructs the asset URL using the development server and webfiles path.
// In production mode, it retrieves the asset entry from the application's entries map.
// If called outside of a server rendering context or if the asset is not found, the function logs a fatal error.
//
// Parameters:
//
//	ctx  - The context, expected to be of type *Context.
//	name - The name of the asset to retrieve.
//
// Returns:
//
//	The URL or path to the requested asset as a string.
func Asset(ctx context.Context, name string) string {
	context, ok := ctx.(*Context)
	if !ok {
		log.Fatal("Asset called outside of server rendering context")
	}
	options := context.app.options

	if options.env == Dev {
		return options.devserver + "/" + options.webfiles + "/" + name
	}
	entry, ok := context.app.entries[name]
	if !ok {
		log.Fatalf("failed to retrieve asset %s", name)
	}
	return entry
}

/*
AsyncNode represents a node that can asynchronously generate an HTML element.
The 'fn' field is a function that returns a pointer to an html.Element, which is
intended to be computed or fetched asynchronously. The 'fallback' field provides
a default html.Node to be used while the asynchronous operation is in progress
or if it fails.
*/
type AsyncNode struct {
	fn       func() *html.Element
	fallback html.Node
}

// Implements the html.Item interface
func (*AsyncNode) Item() {}

// Implements the html.Node interface
func (n *AsyncNode) Render(ctx context.Context, w io.Writer) error {
	context, ok := ctx.(*Context)
	if ok {
		id := util.PrefixedID("pacis")
		context.RegisterChunk(id, n.fn)
		return html.Slot(html.Name(id), html.Fragment(n.fallback)).Render(ctx, w)
	}

	return n.fn().Render(ctx, w)
}

/*
Creates a new AsyncNode

AsyncNode represents a node that can asynchronously generate an HTML element.
The 'fn' field is a function that returns a pointer to an html.Element, which is
intended to be computed or fetched asynchronously. The 'fallback' field provides
a default html.Node to be used while the asynchronous operation is in progress
or if it fails.
*/
func Async(fn func() *html.Element, fallback html.Node) *AsyncNode {
	if fallback == nil {
		fallback = html.Fragment()
	}
	return &AsyncNode{fn: fn, fallback: fallback}
}

/*
SpeculationHook represents a hook used for speculative execution within the system.
It contains the eagerness level for speculation, a reference to the current context,
and a flag indicating whether rendering should occur.
*/
type SpeculationHook struct {
	eagerness SpeculationEagerness
	ctx       *Context
	render    bool
}

// Implements the html.Item interface
func (*SpeculationHook) Item() {}

// Implements the html.Hook interface
func (h *SpeculationHook) Done(el *html.Element) {
	href, ok := el.Attributes["href"]
	if !ok {
		log.Fatal("Speculation hook used in an element without an href attribute")
	}
	h.ctx.RegisterSpeculation(SpeculationRule{
		URLs:      []string{href},
		Eagerness: h.eagerness,
	}, h.render)
}

/*
# Speculation Rules API

	Prefetch returns a speculation hook for registering a speculation rule

SpeculationHook represents a hook used for speculative execution within the system.
It contains the eagerness level for speculation, a reference to the current context,
and a flag indicating whether rendering should occur.

The JSON structure contains one or more fields at the top level, each one representing an action to define speculation rules for. At present the supported actions are:

"prefetch" Optional Experimental
Rules for potential future navigations that should have their associated document response body downloaded, leading to significant performance improvements when those documents are navigated to. Note that none of the subresources referenced by the page are downloaded.

"prerender" Optional Experimental
Rules for potential future navigations that should have their associated documents fully downloaded, rendered, and loaded into an invisible tab. This includes loading all subresources, running all JavaScript, and even loading subresources and performing data fetches started by JavaScript. When those documents are navigated to, navigations will be instant, leading to major performance improvements.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules
*/
func Prefetch(ctx context.Context, eagerness SpeculationEagerness) *SpeculationHook {
	context, ok := ctx.(*Context)
	if !ok {
		log.Fatal("Speculation hook used outside of server rendering context")
	}
	return &SpeculationHook{ctx: context, eagerness: eagerness}
}

/*
# Speculation Rules API

	Prerender returns a speculation hook for registering a speculation rule

SpeculationHook represents a hook used for speculative execution within the system.
It contains the eagerness level for speculation, a reference to the current context,
and a flag indicating whether rendering should occur.

The JSON structure contains one or more fields at the top level, each one representing an action to define speculation rules for. At present the supported actions are:

"prefetch" Optional Experimental
Rules for potential future navigations that should have their associated document response body downloaded, leading to significant performance improvements when those documents are navigated to. Note that none of the subresources referenced by the page are downloaded.

"prerender" Optional Experimental
Rules for potential future navigations that should have their associated documents fully downloaded, rendered, and loaded into an invisible tab. This includes loading all subresources, running all JavaScript, and even loading subresources and performing data fetches started by JavaScript. When those documents are navigated to, navigations will be instant, leading to major performance improvements.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules
*/
func Prerender(ctx context.Context, eagerness SpeculationEagerness) *SpeculationHook {
	context, ok := ctx.(*Context)
	if !ok {
		log.Fatal("Speculation hook used outside of server rendering context")
	}
	return &SpeculationHook{eagerness: eagerness, ctx: context, render: true}
}

// Returns an html.Node for rendering speculation rules script. Should be used inside the <head> tag.
func Speculations(ctx context.Context) html.Node {
	context, ok := ctx.(*Context)
	if !ok {
		return html.Fragment()
	}
	return html.Script(html.Type("speculationrules"), html.JSON(context.specs))
}
