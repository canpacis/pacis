package sheet

import (
	"fmt"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/lucide"
	"github.com/canpacis/pacis/x"
)

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			x.Data("sheet"),
			x.Ref("sheetroot"),
			html.Data("slot", "sheet"),
		)...,
	)
}

func Trigger(items ...html.Item) html.Node {
	return html.Button(
		components.ItemsOf(
			items,

			Open,

			x.Bind("aria-expanded", "opened"),
			x.Bind("data-state", "opened ? 'open' : 'closed'"),

			html.Type("button"),
			html.Data("slot", "sheet-trigger"),
			html.Aria("haspopup", "dialog"),
		)...,
	)
}

func Overlay() html.Node {
	return html.Div(
		x.Show("opened"),
		Close,
		x.Bind("data-state", "opened ? 'open' : 'closed'"),
		html.Data("slot", "sheet-overlay"),
		html.Data("aria-hidden", "true"),
		html.Aria("hidden", "true"),
		html.Attr("x-transition:enter", "transition-opacity duration-150"),
		html.Attr("x-transition:enter-start", "opacity-0"),
		html.Attr("x-transition:enter-end", "opacity-100"),
		html.Attr("x-transition:leave", "transition-opacity duration-150 delay-100"),
		html.Attr("x-transition:leave-start", "opacity-100"),
		html.Attr("x-transition:leave-end", "opacity-0"),
		html.Class("fixed inset-0 z-50 bg-black/50"),
	)
}

type Side = components.Variant

const (
	Top = Side(iota)
	Bottom
	Left
	Right
)

var side = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case Top:
		el.AddClass("inset-x-0 top-0 border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top")
		el.SetAttribute("x-transition:enter-start", "-translate-y-full")
		el.SetAttribute("x-transition:enter-end", "translate-y-0")
		el.SetAttribute("x-transition:leave-start", "translate-y-0")
		el.SetAttribute("x-transition:leave-end", "-translate-y-full")
	case Bottom:
		el.AddClass("inset-x-0 bottom-0 border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom")
		el.SetAttribute("x-transition:enter-start", "translate-y-full")
		el.SetAttribute("x-transition:enter-end", "translate-y-0")
		el.SetAttribute("x-transition:leave-start", "translate-y-0")
		el.SetAttribute("x-transition:leave-end", "translate-y-full")
	case Left:
		el.AddClass("inset-y-0 left-0 h-full w-3/4 border-r data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left sm:max-w-sm")
		el.SetAttribute("x-transition:enter-start", "-translate-x-full")
		el.SetAttribute("x-transition:enter-end", "translate-x-0")
		el.SetAttribute("x-transition:leave-start", "translate-x-0")
		el.SetAttribute("x-transition:leave-end", "-translate-x-full")
	case Right:
		el.AddClass("inset-y-0 right-0 h-full w-3/4  border-l data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right sm:max-w-sm")
		el.SetAttribute("x-transition:enter-start", "translate-x-full")
		el.SetAttribute("x-transition:enter-end", "translate-x-0")
		el.SetAttribute("x-transition:leave-start", "translate-x-0")
		el.SetAttribute("x-transition:leave-end", "translate-x-full")
	default:
		panic(fmt.Sprintf("invalid sheet side variant: %d", v))
	}
})

func Content(items ...html.Item) html.Node {
	items = components.ItemsOf(
		items,
		x.Show("opened"),
		x.Bind("data-state", "opened ? 'open' : 'closed'"),

		html.Data("slot", "sheet-content"),
		html.Attr("x-trap.noscroll", "opened"),
		html.Role("dialog"),
		html.Attr("x-transition:enter", "transition duration-500"),
		html.Attr("x-transition:leave", "transition duration-300"),
		html.Class("fixed flex flex-col z-50 gap-4 bg-background shadow-lg"),
		Right,
		side,
	)
	items = append(items, html.Button(
		Close,
		html.Tabindex("-1"),
		html.Type("button"),
		html.Data("slot", "sheet-close"),
		html.Class("absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground"),

		lucide.X(html.Class("h-4 w-4")),
		html.Span(html.Class("sr-only"), html.Text("Close")),
	))

	return html.Template(
		x.Teleport("body"),

		html.Div(
			Overlay(),
			html.Div(items...),
		),
	)
}

func Header(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "sheet-header"),
			html.Class("flex flex-col space-y-2 text-center sm:text-left p-4"),
		)...,
	)
}

func Footer(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "sheet-footer"),
			html.Class("mt-auto flex flex-col gap-2 p-4"),
		)...,
	)
}

func Title(items ...html.Item) html.Node {
	return html.H2(
		components.ItemsOf(
			items,
			html.Data("slot", "sheet-title"),
			html.Class("text-lg font-semibold text-foreground"),
		)...,
	)
}

func Description(items ...html.Item) html.Node {
	return html.P(
		components.ItemsOf(
			items,
			html.Data("slot", "sheet-description"),
			html.Class("text-sm text-muted-foreground"),
		)...,
	)
}

var Open = OpenOn("click")
var Close = CloseOn("click")

func OpenOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "open($refs.sheetroot)")
}

func CloseOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "close($refs.sheetroot)")
}
