package components

import (
	"fmt"

	h "github.com/canpacis/pacis/ui/html"
)

func Tabs(props ...h.I) h.Element {
	el := h.Div(props...)
	valueattr, _ := el.GetAttribute("value")
	var value any
	if valueattr != nil {
		value = readattr(valueattr)
	}
	id := getid(el)
	el.AddAttribute(X("data", fn("tabs", value, id)))

	return el
}

func TabList(props ...h.I) h.Element {
	return h.Div(
		Join(
			props,
			h.Class("border-b flex gap-0"),
		)...,
	)
}

func TabTrigger(trigger h.Node, props ...h.I) h.Element {
	el := h.Btn(
		Join(props,
			h.Class("relative text-sm h-8 text-muted-foreground font-medium px-4 cursor-pointer after:content-[''] after:w-full after:h-px after:absolute after:left-0 after:-bottom-px after:transition-colors focus-visible:outline-none focus-visible:ring-ring/50 focus-visible:ring-[3px] rounded-t-sm"),
			trigger,
		)...,
	)
	valueattr, ok := el.GetAttribute("value")
	if !ok {
		panic("tab trigger elements need a value attribute")
	}
	value := readattr(valueattr)
	el.AddAttribute(SetTab(value))
	el.AddAttribute(X("bind:class", fmt.Sprintf("value === '%s' && 'after:bg-primary !text-primary'", value)))

	return el
}

func TabContent(props ...h.I) h.Element {
	el := h.Div(
		Join(
			props,
			X("cloak"),
			h.Class("mt-3"),
		)...,
	)
	value, ok := el.GetAttribute("value")
	if !ok {
		panic("tab content elements need a value attribute")
	}
	el.AddAttribute(X("show", fmt.Sprintf("value === '%s'", readattr(value))))

	return el
}

func SetTabOn(event string, value any) h.Attribute {
	return On(event, fn("setTab", value))
}

func SetTab(value any) h.Attribute {
	return SetTabOn("click", value)
}
