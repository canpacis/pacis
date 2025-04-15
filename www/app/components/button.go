package components

import (
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func DocButton(href string, next bool, label Node) Element {
	return Button(
		Href(href),
		Replace(A),
		Class("h-fit min-w-32 justify-start"),
		If(next, Class("ml-auto")),
		ButtonVariantGhost,

		Span(
			Class("flex flex-col gap-px w-full"),
			If(next, Class("items-start")),
			If(!next, Class("items-end")),

			Span(
				Class("text-xs font-light inline"),

				If(next, Text("Up Next")),
				If(!next, Text("Previous")),
			),
			Span(
				Class("text-base font-semibold flex gap-4 items-center w-full"),

				If(!next, icons.ArrowLeft(Class("size-4 mr-auto"))),
				label,
				If(next, icons.ArrowRight(Class("size-4 ml-auto"))),
			),
		),
	)
}
