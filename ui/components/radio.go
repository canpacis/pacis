package components

import (
	"fmt"

	h "github.com/canpacis/pacis/ui/html"
)

func RadioGroup(name h.Attribute, props ...h.I) h.Element {
	el := h.El("fieldset",
		Join(
			props,
			h.Class("space-y-2"),
		)...,
	)

	var value string
	valueattr, hasvalue := el.GetAttribute("value")
	if hasvalue {
		value = readattr(valueattr)
	}

	id := getid(el)

	el.AddAttribute(X("data", fn("radio", readattr(name), value, id)))
	return el
}

func RadioGroupItem(value h.Attribute, props ...h.I) h.Element {
	val := readattr(value)

	return h.Lbl(
		Join(
			props,
			h.Class("flex gap-2 text-sm cursor-pointer items-center"),
			h.Div(
				h.Class("aspect-square h-4 w-4 rounded-full flex justify-center items-center border border-primary text-primary shadow focus:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"),

				h.Span(
					X("show", fmt.Sprintf("value === '%s'", val)),
					h.Class("h-3 w-3 bg-primary rounded-full block"),
				),
				h.Inpt(
					h.Type("radio"),
					h.Class("sr-only"),
					X("bind:name", "name"),
					h.Value(val),
					On("change", fmt.Sprintf("value = '%s'", val)),
				),
			),
		)...,
	)
}
