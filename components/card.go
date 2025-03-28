package components

import r "github.com/canpacis/pacis/renderer"

func CardHeader(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-1.5 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CardTitle(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("leading-none font-semibold"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CardDescription(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("ext-muted-foreground text-sm"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CardAction(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("col-start-2 row-span-2 row-start-1 self-start justify-self-end"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CardContent(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("px-6"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func CardFooter(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("flex items-center px-6 [.border-t]:pt-6"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func Card(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("bg-card text-card-foreground flex flex-col gap-6 rounded-xl border py-6 shadow-sm"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}
