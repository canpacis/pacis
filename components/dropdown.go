package components

import (
	"fmt"

	r "github.com/canpacis/pacis-ui/renderer"
)

func Dropdown(props ...r.I) r.Element {
	props = Join(
		props,
		r.Class("relative"),
		X("data", "dropdown"),
		X("trap.noscroll", "isOpen"),
		On("keydown.esc.window", "closeDropdown"),
	)
	return r.Div(props...)
}

func DropdownTrigger(trigger r.Element) r.Element {
	trigger.AddAttribute(On("click", "isOpen ? closeDropdown() : openDropdown()"))
	trigger.AddAttribute(X("ref", "anchor"))
	return trigger
}

func DropdownContent(props ...r.I) r.Element {
	props = Join(
		props,
		r.Class("min-w-[8rem] overflow-y-auto overflow-x-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"),
		X("cloak"),
		X("show", "isOpen || isKeyboard"),
		X("transition"),
		X("trap", "isKeyboard"),
		Anchor(VBottom, HStart, 8),
		On("click.outside", "closeDropdown"),
		On("keydown.down.prevent", "$focus.wrap().next();"),
		On("keydown.up.prevent", "$focus.wrap().previous();"),
	)
	return r.Div(props...)
}

func DropdownItem(props ...r.I) r.Node {
	props = Join(
		props,
		r.Class("relative flex w-full cursor-default select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&>svg]:size-4 [&>svg]:shrink-0 hover:bg-accent"),
	)
	el := r.Btn(props...)
	ok := true
	// id, ok := el.GetAttribute("id")
	if !ok {
		errset, ok := el.(r.ErrorSetter)
		if ok {
			errset.SetError(fmt.Errorf("dropdown items need a unique id attribute"))
		} else {
			panic("dropdown items need an id attribute")
		}
	} else {
		el.AddAttribute(
			On(
				"click",
				fmt.Sprintf(
					"closeDropdown(), $dispatch('select', { id: '%s' });",
					"TODO PLACE ID HERE",
					// id.GetValue(),
				),
			),
		)
	}

	return r.Try(el, ErrorText)
}
