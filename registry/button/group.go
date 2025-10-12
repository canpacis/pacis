package button

import (
	"log"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"

	"components/ui/separator"
)

type GroupOrientation = components.Variant

const (
	GroupHorizontal = GroupOrientation(iota)
	GroupVertical
)

var orientation = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case GroupHorizontal:
		el.AddClass("[&>*:not(:first-child)]:rounded-l-none [&>*:not(:first-child)]:border-l-0 [&>*:not(:last-child)]:rounded-r-none")
		el.SetAttribute("data-orientation", "horizontal")
	case GroupVertical:
		el.AddClass("flex-col [&>*:not(:first-child)]:rounded-t-none [&>*:not(:first-child)]:border-t-0 [&>*:not(:last-child)]:rounded-b-none")
		el.SetAttribute("data-orientation", "vertical")
	default:
		log.Fatalf("invalid group orientation: %d", v)
	}
})

func Group(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Role("group"),
			html.Data("slot", "button-group"),
			html.Class("flex w-fit items-stretch has-[>[data-slot=button-group]]:gap-2 [&>*]:focus-visible:relative [&>*]:focus-visible:z-10 has-[select[aria-hidden=true]:last-child]:[&>[data-slot=select-trigger]:last-of-type]:rounded-r-md [&>[data-slot=select-trigger]:not([class*='w-'])]:w-fit [&>input]:flex-1"),
			GroupHorizontal,
			orientation,
		)...,
	)
}

func GroupText(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "button-group-text"),
			html.Class("bg-muted shadow-xs flex items-center gap-2 rounded-md border px-4 text-sm font-medium [&_svg:not([class*='size-'])]:size-4 [&_svg]:pointer-events-none"),
		)...,
	)
}

func GroupSeparator(items ...html.Item) html.Node {
	return separator.New(
		components.ItemsOf(
			items,
			html.Data("slot", "button-group-separator"),
			html.Class("bg-input relative !m-0 self-stretch data-[orientation=vertical]:h-auto"),
		)...,
	)
}
