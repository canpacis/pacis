package components

import (
	h "github.com/canpacis/pacis/ui/html"
)

func Dropdown(props ...h.I) h.Element {
	el := h.Div(
		Join(
			props,
			h.Class("relative"),
			X("trap.noscroll", "open"),
			DismissDropdownOn("keydown.esc.window"),
		)...,
	)
	open, hasopen := el.GetAttribute("open")
	_, ok := open.(ComponentAttribute)

	idattr, hasid := el.GetAttribute("id")
	var id string
	if !hasid {
		id = randid()
	} else {
		id = readattr(idattr)
	}
	el.AddAttribute(X("data", fn("dropdown", hasopen && ok, id)))

	return el
}

func DropdownTrigger(trigger h.Element) h.Element {
	trigger.AddAttribute(OpenDropdown)
	trigger.AddAttribute(X("ref", "anchor"))
	return trigger
}

func DropdownContent(props ...h.I) h.Element {
	props = Join(
		props,
		h.Class("min-w-[8rem] overflow-y-auto overflow-x-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 z-50"),
		X("cloak"),
		X("show", "open || usedKeyboard"),
		X("transition"),
		X("trap", "usedKeyboard"),
		Anchor(VBottom, HStart, 8),
		DismissDropdownOn("click.outside"),
		On("keydown.down.prevent", "$focus.wrap().next();"),
		On("keydown.up.prevent", "$focus.wrap().previous();"),
	)
	return h.Div(props...)
}

func DropdownItem(props ...h.I) h.Node {
	props = Join(
		props,
		h.Class("relative flex w-full cursor-default select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&>svg]:size-4 [&>svg]:shrink-0 hover:bg-accent"),
	)
	el := h.Btn(props...)
	idattr, ok := el.GetAttribute("id")
	if !ok {
		panic("dropdown items need an id attribute")
	}
	id := readattr(idattr)

	el.AddAttribute(CloseDropdown(id))
	return el
}

func OpenDropdownOn(event string) h.Attribute {
	return On(event, "openDropdown()")
}

var OpenDropdown = OpenDropdownOn("click")

func CloseDropdownOn(event, value string) h.Attribute {
	return On(event, fn("closeDropdown", value))
}

func CloseDropdown(value string) h.Attribute {
	return CloseDropdownOn("click", value)
}

func DismissDropdownOn(event string) h.Attribute {
	return On(event, "dismissDropdown()")
}

var DismissDropdown = DismissDropdownOn("click")
