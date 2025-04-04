package components

import h "github.com/canpacis/pacis/ui/html"

func CardHeader(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-1.5 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

func CardTitle(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("leading-none font-semibold"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

func CardDescription(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("ext-muted-foreground text-sm"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

func CardAction(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("col-start-2 row-span-2 row-start-1 self-start justify-self-end"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

func CardContent(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("px-6"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

func CardFooter(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("flex items-center px-6 [.border-t]:pt-6"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

func Card(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("bg-card text-card-foreground flex flex-col gap-6 rounded-xl border py-6 shadow-sm"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}
