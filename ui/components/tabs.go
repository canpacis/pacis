package components

import (
	"bytes"
	"context"
	"fmt"

	h "github.com/canpacis/pacis/ui/html"
)

func Tabs(props ...h.I) h.Element {
	el := h.Div(Join(props, X("data", "tabs"))...)
	value, ok := el.GetAttribute("value")
	if ok {
		var buf bytes.Buffer
		value.Render(context.Background(), &buf)

		el.AddAttribute(X("init", fmt.Sprintf("value = '%s'", buf.String())))
	}

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
			h.Class("relative text-sm h-8 text-muted-foreground font-medium px-4 cursor-pointer after:content-[''] after:w-full after:h-px after:absolute after:left-0 after:-bottom-px after:transition-colors"),
			trigger,
		)...,
	)
	value, ok := el.GetAttribute("value")
	if !ok {
		panic("tab trigger elements need a value attribute")
	}
	var buf bytes.Buffer
	value.Render(context.Background(), &buf)
	el.AddAttribute(SetActiveTab(buf.String()))
	el.AddAttribute(X("bind:class", fmt.Sprintf("value === '%s' && 'after:bg-primary !text-primary'", buf.String())))

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
	var buf bytes.Buffer
	value.Render(context.Background(), &buf)
	el.AddAttribute(X("show", fmt.Sprintf("value === '%s'", buf.String())))

	return el
}

func SetActiveTab(value string) h.Attribute {
	return On("click", fmt.Sprintf("setActiveTab('%s')", value))
}
