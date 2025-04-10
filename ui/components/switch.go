package components

import h "github.com/canpacis/pacis/ui/html"

func Switch(props ...h.I) h.Element {
	el := h.S(props...)
	_, checked := el.GetAttribute("checked")
	idattr, hasid := el.GetAttribute("id")

	var id string
	if !hasid {
		id = randid()
	} else {
		id = readattr(idattr)
	}

	return h.Lbl(
		Join(
			props,
			h.Class("flex gap-2 text-sm cursor-pointer"),
			X("data", fn("switch_", checked, id)),
			h.Div(
				h.Data(":state", "checked ? 'checked' : 'unchecked'"),
				h.Class("peer inline-flex h-5 w-9 shrink-0 items-center rounded-full border-2 border-transparent shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-primary data-[state=unchecked]:bg-input"),
				h.Span(
					h.Data(":state", "checked ? 'checked' : 'unchecked'"),
					h.Class("pointer-events-none block h-4 w-4 rounded-full bg-background shadow-lg ring-0 transition-transform data-[state=checked]:translate-x-4 data-[state=unchecked]:translate-x-0"),
				),
				h.Inpt(
					h.ID(id),
					h.Type("checkbox"),
					h.Class("sr-only"),
					X("bind:checked", "checked"),
					ToggleSwitchOn("change"),
				),
			),
		)...,
	)
}

// Returns an attribute that toggles the related switch upon given event
func ToggleSwitchOn(event string) h.Attribute {
	return On(event, "toggleSwitch()")
}

// An attribute that toggles the related switch on click
var ToggleSwitch = ToggleSwitchOn("click")
