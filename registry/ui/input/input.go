package input

import (
	"fmt"

	"components/ui/button"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

func New(items ...html.Item) html.Node {
	return html.Input(
		components.ItemsOf(
			items,
			html.Data("slot", "input"),
			html.Class("file:text-foreground placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 border-input h-9 w-full min-w-0 rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"),
		)...,
	)
}

func Group(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "input-group"),
			html.Role("group"),
			html.Class("group/input-group border-input dark:bg-input/30 shadow-xs relative flex w-full items-center rounded-md border outline-none transition-[color,box-shadow] h-9 has-[>textarea]:h-auto has-[>[data-align=inline-start]]:[&>input]:pl-2 has-[>[data-align=inline-end]]:[&>input]:pr-2 has-[>[data-align=block-start]]:h-auto has-[>[data-align=block-start]]:flex-col has-[>[data-align=block-start]]:[&>input]:pb-3 has-[>[data-align=block-end]]:h-auto has-[>[data-align=block-end]]:flex-col has-[>[data-align=block-end]]:[&>input]:pt-3 has-[[data-slot=input-group-control]:focus-visible]:ring-ring has-[[data-slot=input-group-control]:focus-visible]:ring-1 has-[[data-slot][aria-invalid=true]]:ring-destructive/20 has-[[data-slot][aria-invalid=true]]:border-destructive dark:has-[[data-slot][aria-invalid=true]]:ring-destructive/40"),
		)...,
	)
}

type AddonAlign = components.Variant

const (
	AddonInlineStart = AddonAlign(iota)
	AddonInlineEnd
	AddonBlockStart
	AddonBlockEnd
)

var addon = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case AddonInlineStart:
		el.AddClass("order-first pl-3 has-[>button]:ml-[-0.45rem] has-[>kbd]:ml-[-0.35rem]")
		el.SetAttribute("data-align", "inline-start")
	case AddonInlineEnd:
		el.AddClass("order-last pr-3 has-[>button]:mr-[-0.4rem] has-[>kbd]:mr-[-0.35rem]")
		el.SetAttribute("data-align", "inline-end")
	case AddonBlockStart:
		el.AddClass("[.border-b]:pb-3 order-first w-full justify-start px-3 pt-3 group-has-[>input]/input-group:pt-2.5")
		el.SetAttribute("data-align", "block-start")
	case AddonBlockEnd:
		el.AddClass("[.border-t]:pt-3 order-last w-full justify-start px-3 pb-3 group-has-[>input]/input-group:pb-2.5")
		el.SetAttribute("data-align", "block-end")
	default:
		panic(fmt.Sprintf("invalid input group addon align variant: %d", v))
	}
})

func GroupAddon(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "input-group-addon"),
			html.Role("group"),
			addon,
		)...,
	)
}

type GroupButtonSize = components.Size

const (
	ButtonXs = GroupButtonSize(iota)
	ButtonSm
	ButtonIconXs
	ButtonIconSm
)

var size = components.NewSizeApplier(func(el *html.Element, s components.Size) {
	switch s {
	case ButtonXs:
		el.AddClass("h-6 gap-1 rounded-[calc(var(--radius)-5px)] px-2 has-[>svg]:px-2 [&>svg:not([class*='size-'])]:size-3.5")
		el.SetAttribute("data-size", "xs")
	case ButtonSm:
		el.AddClass("h-8 gap-1.5 rounded-md px-2.5 has-[>svg]:px-2.5")
		el.SetAttribute("data-size", "sm")
	case ButtonIconXs:
		el.AddClass("size-6 rounded-[calc(var(--radius)-5px)] p-0 has-[>svg]:p-0")
		el.SetAttribute("data-size", "icon-xs")
	case ButtonIconSm:
		el.AddClass("size-8 p-0 has-[>svg]:p-0")
		el.SetAttribute("data-size", "icon-sm")
	default:
		panic(fmt.Sprintf("invalid group button size: %d", s))
	}
})

func GroupButton(items ...html.Item) html.Node {
	return button.New(components.ItemsOf(items, size)...)
}

func GroupText(items ...html.Item) html.Node {
	return html.Span(components.ItemsOf(items, html.Class("text-muted-foreground flex items-center gap-2 text-sm [&_svg:not([class*='size-'])]:size-4 [&_svg]:pointer-events-none"))...)
}

func GroupInput(items ...html.Item) html.Node {
	return New(components.ItemsOf(items, html.Data("slot", "input-group-control"), html.Class("flex-1 rounded-none border-0 bg-transparent shadow-none focus-visible:ring-0 dark:bg-transparent"))...)
}

func GroupTextarea(items ...html.Item) html.Node {
	return html.Textarea(
		components.ItemsOf(
			items,
			html.Data("slot", "input-group-control"),
			html.Class("flex-1 resize-none rounded-none border-0 bg-transparent py-3 shadow-none focus-visible:ring-0 dark:bg-transparent"),
		)...,
	)
}
