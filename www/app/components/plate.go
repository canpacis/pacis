package components

import (
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
)

func Plate(node Node, props ...I) Element {
	return Div(
		Join(
			props,
			Class("flex justify-center items-center gap-4 min-h-40 w-full border rounded-lg p-4 md:p-16 my-4"),

			node,
		)...,
	)
}
