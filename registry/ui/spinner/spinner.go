package spinner

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/lucide"
)

func New(items ...html.Item) html.Node {
	return lucide.LoaderCircle(
		components.ItemsOf(
			items,
			html.Role("status"),
			html.Aria("label", "loading"),
			html.Class("size-4 animate-spin"),
		)...,
	)
}
