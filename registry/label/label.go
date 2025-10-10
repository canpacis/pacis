package label

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

func New(items ...html.Item) html.Node {
	return html.Label(
		components.ItemsOf(
			items,
			html.Data("slot", "label"),
			html.Class("text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"),
		)...,
	)
}
