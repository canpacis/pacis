package pagination

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/lucide"
)

func New(items ...html.Item) html.Node {
	return html.Nav(
		components.ItemsOf(
			items,
			html.Role("navigation"),
			html.Aria("label", "pagination"),
			html.Class("mx-auto flex w-full justify-center"),
		)...,
	)
}

func Content(items ...html.Item) html.Node {
	return html.Ul(
		components.ItemsOf(
			items,
			html.Class("flex flex-row items-center gap-1"),
		)...,
	)
}

func Item(items ...html.Item) html.Node {
	return html.Li(items...)
}

func Link(active bool, items ...html.Item) html.Node {
	return html.A(
		components.ItemsOf(
			items,
			html.If(active, html.Aria("current", "")),
			html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0"),
			html.If(active, html.Class("border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground")),
			html.If(!active, html.Class("hover:bg-accent hover:text-accent-foreground")),
		)...,
	)
}

func Previous(items ...html.Item) html.Node {
	return Link(false,
		components.ItemsOf(
			items,
			html.Aria("label", "Go to previous page"),
			html.Class("gap-1 pl-2.5"),

			lucide.ChevronLeft(html.Class("h-4 w-4")),
			html.Span(html.Text("Previous")),
		)...,
	)
}

func Next(items ...html.Item) html.Node {
	return Link(false,
		components.ItemsOf(
			items,
			html.Aria("label", "Go to next page"),
			html.Class("gap-1 pl-2.5"),

			html.Span(html.Text("Next")),
			lucide.ChevronRight(html.Class("h-4 w-4")),
		)...,
	)
}

func Ellipsis(items ...html.Item) html.Node {
	return html.Span(
		components.ItemsOf(
			items,
			html.Aria("hidden", ""),
			html.Class("flex h-9 w-9 items-center justify-center"),

			lucide.Ellipsis(html.Class("h-4 w-4")),
			html.Span(html.Class("sr-only"), html.Text("More pages")),
		)...,
	)
}
