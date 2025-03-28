package components

import r "github.com/canpacis/pacis/renderer"

func Collapsible(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		D{"open": false},
		r.Data(":state", "open ? 'opened' : 'closed'"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CollapsibleTrigger(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		On("click", "open = !open"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CollapsibleContent(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Attr("x-show", "open"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}
