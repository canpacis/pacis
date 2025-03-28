package components

import r "github.com/canpacis/pacis/renderer"

func Dropdown(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("relative "),
		D{
			"open":         false,
			"usedKeyboard": false,
		},
		On("keydown.esc.window", "open = false, usedKeyboard = false"),
	}

	ps = append(ps, props...)
	return r.Div(ps...)
}

func DropdownTrigger(trigger *r.Element) *r.Element {
	return r.Clone(trigger, On("click", "open = !open"))
}

func DropdownContent(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("absolute top-[calc(100%+0.325rem)] z-50 min-w-[8rem] overflow-y-auto overflow-x-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"),
		r.Attr("x-cloak"),
		r.Attr("x-show", "open || usedKeyboard"),
		r.Attr("x-transition"),
		r.Attr("x-trap", "usedKeyboard"),
		On("click.outside", "open = false, usedKeyboard = false"),
		On("keydown.down.prevent", "$focus.wrap().next()"),
		On("keydown.up.prevent", "$focus.wrap().previous()"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func DropdownItem(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("relative flex w-full cursor-default select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&>svg]:size-4 [&>svg]:shrink-0 hover:bg-accent"),
		On("click", "open = false, usedKeyboard = false"),
	}
	ps = append(ps, props...)

	return r.Btn(ps...)
}
