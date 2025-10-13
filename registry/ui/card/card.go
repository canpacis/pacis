package card

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "card"),
			html.Class("rounded-xl border bg-card text-card-foreground shadow"),
		)...,
	)
}

func Header(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "card-header"),
			html.Class("flex flex-col space-y-1.5 p-6"),
		)...,
	)
}

func Title(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "card-title"),
			html.Class("font-semibold leading-none tracking-tight"),
		)...,
	)
}

func Description(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "card-description"),
			html.Class("text-sm text-muted-foreground"),
		)...,
	)
}

func Content(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "card-content"),
			html.Class("p-6 pt-0"),
		)...,
	)
}

func Footer(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "card-footer"),
			html.Class("flex items-center p-6 pt-0"),
		)...,
	)
}
