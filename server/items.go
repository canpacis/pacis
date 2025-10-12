package server

import (
	"context"
	"log"
	"net/http"

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
func Asset(app *App, name string) string {
	if app.options.env == Dev {
		return app.options.devserver + "/" + app.options.webfiles + "/" + name
	}
	entry, ok := app.entries[name]
	if !ok {
		log.Fatalf("failed to retrieve asset %s", name)
	}
	return entry
}

func Async(comp html.Component, fallback html.Node) html.Node {
	id := util.PrefixedID("pacis")
	if fallback == nil {
		fallback = html.Fragment()
	}

	return html.Component(func(ctx context.Context) html.Node {
		serverctx, ok := ctx.(*serverctx)
		if ok {
			serverctx.async = append(serverctx.async, async{
				id:   id,
				comp: comp,
			})
		}
		return html.Slot(html.Name(id), fallback)
	})
}

func Redirect(ctx context.Context, to string) html.Node {
	context, ok := ctx.(*serverctx)
	if !ok {
		log.Fatal("Redirect node used outside of server rendering context")
	}
	context.redirect = &redirect{status: http.StatusFound, to: to}
	return html.Fragment()
}

func RedirectWith(ctx context.Context, to string, status int) html.Node {
	context, ok := ctx.(*serverctx)
	if !ok {
		log.Fatal("RedirectWith node used outside of server rendering context")
	}
	context.redirect = &redirect{status: status, to: to}
	return html.Fragment()
}

/*
# Speculation Rules API

A string providing a hint to the browser as to how eagerly it should prefetch/prerender
link targets in order to balance performance advantages against resource overheads.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#eagerness
*/
type SpeculationEagerness string

const (
	/*
		# Speculation Rules API

		The author thinks the link is very likely to be followed, and/or the document may take
		significant time to fetch. Prefetch/prerender should start as soon as possible, subject
		only to considerations such as user preferences and resource limits.

		https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#immediate
	*/
	ImmediateEagerness = SpeculationEagerness("immediate")
	/*
		# Speculation Rules API

		The author wants to prefetch/prerender a large number of navigations, as early as possible.
		Prefetch/prerender should start on any slight suggestion that a link may be followed.
		For example, the user could move their mouse cursor towards the link, hover/focus it for a
		moment, or pause scrolling with the link in a prominent place.

		https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#eager
	*/
	EagerEagerness = SpeculationEagerness("eager")
	/*
		# Speculation Rules API

		The author is looking for a balance between eager and conservative. Prefetch/prerender
		should start when there is a reasonable suggestion that the user will follow a link
		in the near future. For example, the user could scroll a link into the viewport and
		hover/focus it for some time.

		https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#moderate
	*/
	ModerateEagerness = SpeculationEagerness("moderate")
	/*
		# Speculation Rules API

		The author wishes to get some benefit from speculative loading with a fairly small
		tradeoff of resources. Prefetch/prerender should start only when the user is starting
		to click on the link, for example on mousedown or pointerdown.

		https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#conservative
	*/
	ConservativeEagerness = SpeculationEagerness("conservative")
)

// /*
// # Speculation Rules API

// The JSON structure contains one or more fields at the top level, each one representing an
// action to define speculation rules for. At present the supported actions are:

// "prefetch" Optional Experimental
// Rules for potential future navigations that should have their associated document response
// body downloaded, leading to significant performance improvements when those documents are
// navigated to. Note that none of the subresources referenced by the page are downloaded.

// "prerender" Optional Experimental
// Rules for potential future navigations that should have their associated documents fully
// downloaded, rendered, and loaded into an invisible tab. This includes loading all subresources,
// running all JavaScript, and even loading subresources and performing data fetches started by
// JavaScript. When those documents are navigated to, navigations will be instant, leading to
// major performance improvements.

// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules
// */

/*
SpeculationProperty represents a hook used for speculative execution within the system.
It contains the eagerness level for speculation, a reference to the current context,
and a flag indicating whether rendering should occur.
*/
type SpeculationProperty struct {
	eagerness SpeculationEagerness
	render    bool
}

func (*SpeculationProperty) LifeCycle() html.PropertyLifeCycle {
	return html.LifeCycleStatic
}

// Implements the html.Item interface
func (*SpeculationProperty) Item() {}

// Implements the html.Hook interface
func (h *SpeculationProperty) Apply(ctx context.Context, el *html.Element) {
	context, ok := ctx.(*serverctx)
	if !ok {
		log.Fatal("Speculation property used outside of server rendering context")
	}
	href := el.GetAttribute("href")
	if len(href) == 0 {
		log.Fatal("Speculation property used in an element without an href attribute")
	}
	if h.render {
		context.specs.Prerender = append(context.specs.Prerender, specrule{
			URLs:      []string{href},
			Eagerness: h.eagerness,
		})
	} else {
		context.specs.Prefetch = append(context.specs.Prefetch, specrule{
			URLs:      []string{href},
			Eagerness: h.eagerness,
		})
	}
}

/*
# Speculation Rules API

	Prefetch returns a speculation hook for registering a speculation rule

SpeculationHook represents a hook used for speculative execution within the system.
It contains the eagerness level for speculation, a reference to the current context,
and a flag indicating whether rendering should occur.

The JSON structure contains one or more fields at the top level, each one representing
an action to define speculation rules for. At present the supported actions are:

"prefetch" Optional Experimental
Rules for potential future navigations that should have their associated document response
body downloaded, leading to significant performance improvements when those documents are
navigated to. Note that none of the subresources referenced by the page are downloaded.

"prerender" Optional Experimental
Rules for potential future navigations that should have their associated documents fully
downloaded, rendered, and loaded into an invisible tab. This includes loading all subresources,
running all JavaScript, and even loading subresources and performing data fetches started by
JavaScript. When those documents are navigated to, navigations will be instant, leading to
major performance improvements.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules
*/
// func Prefetch(ctx context.Context, eagerness SpeculationEagerness) *SpeculationHook {
// 	context, ok := ctx.(*Page)
// 	if !ok {
// 		log.Fatal("Speculation hook used outside of server rendering context")
// 	}
// 	return &SpeculationHook{ctx: context, eagerness: eagerness}
// }

// /*
// # Speculation Rules API

// 	Prerender returns a speculation hook for registering a speculation rule

// SpeculationHook represents a hook used for speculative execution within the system.
// It contains the eagerness level for speculation, a reference to the current context,
// and a flag indicating whether rendering should occur.

// The JSON structure contains one or more fields at the top level, each one representing an
// action to define speculation rules for. At present the supported actions are:

// "prefetch" Optional Experimental
// Rules for potential future navigations that should have their associated document response
// body downloaded, leading to significant performance improvements when those documents are
// navigated to. Note that none of the subresources referenced by the page are downloaded.

// "prerender" Optional Experimental
// Rules for potential future navigations that should have their associated documents fully
// downloaded, rendered, and loaded into an invisible tab. This includes loading all subresources,
// running all JavaScript, and even loading subresources and performing data fetches started by
// JavaScript. When those documents are navigated to, navigations will be instant, leading to
// major performance improvements.

// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules
// */
// func Prerender(ctx context.Context, eagerness SpeculationEagerness) *SpeculationHook {
// 	context, ok := ctx.(*Page)
// 	if !ok {
// 		log.Fatal("Speculation hook used outside of server rendering context")
// 	}
// 	return &SpeculationHook{eagerness: eagerness, ctx: context, render: true}
// }

func Prerender(eagerness SpeculationEagerness) html.Property {
	return &SpeculationProperty{eagerness: eagerness, render: true}
}

func Speculations(ctx context.Context) html.Node {
	context, ok := ctx.(*serverctx)
	if !ok {
		return html.Fragment()
	}
	if context.specs.isempty() {
		return html.Fragment()
	}
	return html.Script(html.Type("speculationrules"), html.JSON(context.specs))
}
