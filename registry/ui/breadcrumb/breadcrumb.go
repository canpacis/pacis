package breadcrumb

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/lucide"
)

func New(items ...html.Item) html.Node {
	return html.Nav(
		components.ItemsOf(
			items,
			html.Data("slot", "breadcrumb"),
			html.Aria("label", "breadcrumb"),
		)...,
	)
}

func List(items ...html.Item) html.Node {
	return html.Ol(
		components.ItemsOf(
			items,
			html.Data("slot", "breadcrumb-list"),
			html.Class("flex flex-wrap items-center gap-1.5 break-words text-sm text-muted-foreground sm:gap-2.5"),
		)...,
	)
}

func Item(items ...html.Item) html.Node {
	return html.Li(
		components.ItemsOf(
			items,
			html.Data("slot", "breadcrumb-item"),
			html.Class("inline-flex items-center gap-1.5"),
		)...,
	)
}

func Link(items ...html.Item) html.Node {
	return html.A(
		components.ItemsOf(
			items,
			html.Data("slot", "breadcrumb-link"),
			html.Class("transition-colors hover:text-foreground"),
		)...,
	)
}

func Page(items ...html.Item) html.Node {
	return html.Span(
		components.ItemsOf(
			items,
			html.Role("link"),
			html.Aria("disabled", "true"),
			html.Aria("current", "page"),
			html.Data("slot", "breadcrumb-page"),
			html.Class("font-normal text-foreground"),
		)...,
	)
}

func Separator(items ...html.Item) html.Node {
	return html.Li(
		components.ItemsOf(
			items,
			html.Role("presentation"),
			html.Aria("hidden", "true"),
			html.Data("slot", "breadcrumb-separator"),
			html.Class("[&>svg]:w-3.5 [&>svg]:h-3.5"),
		)...,
	)
}

func Ellipsis(items ...html.Item) html.Node {
	return html.Span(
		components.ItemsOf(
			items,
			html.Role("presentation"),
			html.Aria("hidden", "true"),
			html.Data("slot", "breadcrumb-ellipsis"),
			html.Class("flex h-9 w-9 items-center justify-center"),

			lucide.Ellipsis(html.Class("h-4 w-4")),
			html.Span(html.Class("sr-only"), html.Text("More")),
		)...,
	)
}
