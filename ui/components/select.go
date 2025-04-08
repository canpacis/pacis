package components

import (
	"bytes"
	"context"
	"fmt"

	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func Select(props ...h.I) h.Element {
	el := h.Div(
		Join(
			props,
			X("data", "select"), h.Class("relative"),
			X("trap.noscroll", "isOpen"),
			On("keydown.esc.window", "closeDropdown"),
		)...,
	)
	name, ok := el.GetAttribute("name")

	if ok {
		el.AddNode(Input(name, X("bind:value", "value"), h.Class("sr-only")))
	} else {
		panic("select element requires a name attribute")
	}

	_, clearable := el.GetAttribute("clearable")
	if clearable {
		el.AddAttribute(X("init", "clearable = true"))
	}

	return el
}

func SelectTrigger(trigger h.Element, selected h.Element) h.Element {
	selected.AddAttribute(X("show", "value !== null"))
	trigger.AddAttribute(X("show", "value === null"))

	return h.Div(
		h.Class("relative"),

		h.Btn(
			h.Class("flex flex-1 h-9 w-full items-center whitespace-nowrap rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm ring-offset-background data-[placeholder]:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1"),
			X("ref", "anchor"),
			On("click", "openSelect()"),

			trigger,
			selected,
		),
		h.Span(
			h.Class("absolute top-0 bottom-0 right-2 my-auto flex gap-2 justify-center items-center pointer-events-none"),

			Button(
				ButtonSizeIcon,
				ButtonVariantGhost,
				h.Class("w-5 h-5 pointer-events-auto"),
				X("show", "value !== null && clearable"),
				On("click", "value = null"),

				icons.X(h.Class("size-4 opacity-50")),
			),
			icons.ChevronDown(h.Class("size-4 opacity-50")),
		),
	)
}

func SelectContent(props ...h.I) h.Element {
	return h.Div(
		Join(
			props,
			h.Class("z-50 min-w-[8rem] overflow-y-auto overflow-x-hidden rounded-md border bg-popover text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 origin-[--radix-select-content-transform-origin] data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1 p-1 w-full"),
			h.Data(":state", "isOpen ? 'open' : 'closed'"),
			X("cloak"),
			X("show", "isOpen || isKeyboard"),
			X("transition"),
			X("trap", "isKeyboard"),
			Anchor(VBottom, HStart, 8),
			On("click.outside", "closeSelect(null, true)"),
			On("keydown.down.prevent", "$focus.wrap().next();"),
			On("keydown.up.prevent", "$focus.wrap().previous();"),
		)...,
	)
}

func SelectItem(props ...h.I) h.Element {
	el := h.Btn(
		Join(
			props,
			h.Class("relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-2 pr-8 text-sm outline-none focus:bg-accent hover:bg-accent focus:text-accent-foreground hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"),
		)...,
	)
	value, ok := el.GetAttribute("value")
	if !ok {
		panic("select item elements need a value attribute")
	}
	var buf bytes.Buffer
	value.Render(context.Background(), &buf)
	el.AddAttribute(On("click", fmt.Sprintf("closeSelect('%s', false)", buf.String())))

	return el
}

func SelectSeperator() h.Element {
	return h.Span(h.Class("-mx-1 my-1 h-px bg-muted block"))
}

func SelectLabel(label string, props ...h.I) h.Element {
	return h.Span(Join(props, h.Class("px-2 py-1.5 text-xs font-semibold text-muted-foreground"), h.Text(label))...)
}
