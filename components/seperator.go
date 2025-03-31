package components

import r "github.com/canpacis/pacis-ui/renderer"

func Seperator(orientation Orientation, props ...r.I) r.Element {
	return r.Div(
		Join(
			props,
			r.Role("none"),
			r.Data("orientation", orientation.String()),
			r.Class("bg-border"),
			r.If(orientation == OHorizontal, r.Class("h-[1px] w-full")),
			r.If(orientation == OVertical, r.Class("h-full w-[1px]")),
		)...,
	)
}
