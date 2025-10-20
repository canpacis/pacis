package separator

import (
	"fmt"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

type Orientation = components.Variant

const (
	Horizontal = Orientation(iota)
	Vertical
)

var orientation = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case Horizontal:
		el.AddClass("h-[1px] w-full")
		el.SetAttribute("data-orientation", "horizontal")
	case Vertical:
		el.AddClass("h-full w-[1px]")
		el.SetAttribute("data-orientation", "vertical")
	default:
		panic(fmt.Sprintf("invalid group orientation: %d", v))
	}
})

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Role("none"),
			html.Data("slot", "separator"),
			html.Class("shrink-0 bg-border"),
			Horizontal,
			orientation,
		)...,
	)
}
