package components

import r "github.com/canpacis/pacis-ui/renderer"

func Collapsible(props ...r.I) r.Element {
	ps := []r.I{
		D{"open": false},
		r.Data(":state", "open ? 'opened' : 'closed'"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CollapsibleTrigger(trigger r.Element) r.Element {
	trigger.AddAttribute(On("click", "open = !open;"))
	return trigger
}

func CollapsibleContent(content r.Element) r.Element {
	content.AddAttribute(r.Attr("x-show", "open"))
	return content
}
