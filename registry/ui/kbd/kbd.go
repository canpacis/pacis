package kbd

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

func New(items ...html.Item) html.Node {
	return html.Kbd(
		components.ItemsOf(
			items,
			html.Data("slot", "kbd"),
			html.Class("bg-muted text-muted-foreground pointer-events-none inline-flex h-5 w-fit min-w-5 select-none items-center justify-center gap-1 rounded-sm px-1 font-sans text-xs font-medium [&_svg:not([class*='size-'])]:size-3 [[data-slot=tooltip-content]_&]:bg-background/20 [[data-slot=tooltip-content]_&]:text-background dark:[[data-slot=tooltip-content]_&]:bg-background/10"),
		)...,
	)
}

func Group(items ...html.Item) html.Node {
	return html.Kbd(
		components.ItemsOf(
			items,
			html.Data("slot", "kbd-group"),
			html.Class("inline-flex items-center gap-1"),
		)...,
	)
}
