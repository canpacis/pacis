package components

import (
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
)

func SectionTitle(text Node) Element {
	return Div(
		Class("my-4"),

		H2(
			Class("scroll-m-20 text-xl font-bold tracking-tight"),

			text,
		),
		Seperator(OHorizontal),
	)
}
