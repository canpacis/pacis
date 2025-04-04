package components

import h "github.com/canpacis/pacis/ui/html"

func Seperator(orientation Orientation, props ...h.I) h.Element {
	return h.Div(
		Join(
			props,
			h.Role("none"),
			h.Data("orientation", orientation.String()),
			h.Class("bg-border my-2"),
			h.If(orientation == OHorizontal, h.Class("h-[1px] w-full")),
			h.If(orientation == OVertical, h.Class("h-full w-[1px]")),
		)...,
	)
}
