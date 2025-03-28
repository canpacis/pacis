package components

import r "github.com/canpacis/pacis/renderer"

func Collapsible(props ...r.I) r.Element {
	ps := []r.I{
		D{"open": false},
		r.Data(":state", "open ? 'opened' : 'closed'"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CollapsibleTrigger(props ...r.I) r.Element {
	ps := []r.I{
		On("click", "open = !open"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CollapsibleContent(props ...r.I) r.Element {
	ps := []r.I{
		r.Attr("x-show", "open"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}
