package components

import (
	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func Input(props ...h.I) h.Element {
	input := h.Inpt(
		Join(
			props,
			h.Class("flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-base shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm appearance-none"),
		)...,
	)
	typ, ok := input.GetAttribute("type")

	return h.Div(
		h.Class("relative h-fit"),

		D{"input": ""},
		X("init", "input = $el.querySelector('input')"),

		input,
		h.IfFn(ok, func() h.Renderer {
			// input[type=number]
			return h.If(readattr(typ) == "number", h.Span(
				h.Class("absolute top-0 right-0 h-full flex flex-col text-muted-foreground"),

				h.Btn(
					h.Type("button"),
					h.Class("bg-transparent border-l border-b flex items-center justify-center px-1 flex-1 hover:bg-accent/50"),
					On("click", "input.value = Number(input.value) + 1"),

					icons.ChevronUp(h.Class("size-3.5 pointer-events-none")),
				),
				h.Btn(
					h.Type("button"),
					h.Class("bg-transparent border-l flex items-center justify-center px-1 flex-1 hover:bg-accent/50"),
					On("click", "input.value = Number(input.value) - 1"),

					icons.ChevronDown(h.Class("size-3.5 pointer-events-none")),
				),
			))
		}),
	)
}
