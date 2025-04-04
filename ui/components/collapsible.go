package components

import h "github.com/canpacis/pacis/ui/html"

func Collapsible(props ...h.I) h.Element {
	ps := []h.I{
		D{"open": false},
		h.Data(":state", "open ? 'opened' : 'closed'"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

func CollapsibleTrigger(trigger h.Element) h.Element {
	trigger.AddAttribute(On("click", "open = !open;"))
	return trigger
}

func CollapsibleContent(content h.Element) h.Element {
	content.AddAttribute(h.Attr("x-show", "open"))
	return content
}
