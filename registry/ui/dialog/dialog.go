package dialog

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/lucide"
	"github.com/canpacis/pacis/x"
)

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			x.Data("dialog"),
			x.Ref("dialogroot"),
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
			html.Aria("has-popup", "dialog"),
		)...,
	)
}

func Overlay() html.Node {
	return html.Div(
		x.Show("opened"),
		Close,
		x.Bind("data-state", "opened ? 'open' : 'closed'"),
		html.Data("slot", "dialog-overlay"),
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

func Content(items ...html.Item) html.Node {
	items = components.ItemsOf(
		items,
		x.Show("opened"),
		x.Bind("data-state", "opened ? 'open' : 'closed'"),
		html.Data("slot", "dialog-cotent"),
		html.Role("dialog"),
		html.Attr("x-trap.noscroll", "opened"),
		html.Attr("x-transition:enter", "transition-opacity duration-150 delay-100"),
		html.Attr("x-transition:enter-start", "opacity-0"),
		html.Attr("x-transition:enter-end", "opacity-100"),
		html.Attr("x-transition:leave", "transition-opacity duration-150"),
		html.Attr("x-transition:leave-start", "opacity-100"),
		html.Attr("x-transition:leave-end", "opacity-0"),
		html.Class("fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg sm:rounded-lg"),
	)

	items = append(items, html.Button(
		Close,
		html.Tabindex("-1"),
		html.Type("button"),
		html.Data("slot", "dialog-close"),
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
	return html.Div(components.ItemsOf(items, html.Class("flex flex-col space-y-1.5 text-center sm:text-left"))...)
}

func Footer(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Class("flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2"))...)
}

func Title(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "dialog-title"), html.Class("text-lg font-semibold leading-none tracking-tight"))...)
}

func Description(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "dialog-description"), html.Class("text-sm text-muted-foreground"))...)
}

var Open = OpenOn("click")
var Close = CloseOn("click")

func OpenOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "open($refs.dialogroot)")
}

func CloseOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "close($refs.dialogroot)")
}
