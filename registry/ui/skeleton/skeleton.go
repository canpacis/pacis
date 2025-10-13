package skeleton

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "skeleton"),
			html.Class("animate-pulse rounded-md bg-primary/10"),
		)...,
	)
}
