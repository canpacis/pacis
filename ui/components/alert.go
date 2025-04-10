package components

import (
	h "github.com/canpacis/pacis/ui/html"
)

// Title slot for the alert component
func AlertTitle(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("col-start-2 line-clamp-1 min-h-4 font-medium tracking-tight"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

// Description slot for the alert component
func AlertDescription(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("text-muted-foreground col-start-2 grid justify-items-start gap-1 text-sm [&_p]:leading-relaxed"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

var (
	AlertVariantDefault     = &GroupedClass{"alert", "bg-card text-card-foreground", true}
	AlertVariantDestructive = &GroupedClass{"alert", "text-destructive bg-card [&>svg]:text-current *:data-[slot=alert-description]:text-destructive/90", false}
)

/*
	Displays a callout for user attention.

Usage:

	Alert(
		icons.Code(),
		AlertTitle(Text("Heads up!")),
		AlertDescription(Text("You can us Go tho create great UI's")),
	)
*/
func Alert(props ...h.I) h.Element {
	return h.Div(
		Join(
			props,
			AlertVariantDefault,
			h.Class("relative w-full rounded-lg border px-4 py-3 text-sm grid has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr] grid-cols-[0_1fr] has-[>svg]:gap-x-3 gap-y-0.5 items-start [&>svg]:size-4 [&>svg]:translate-y-0.5 [&>svg]:text-current"),
			h.Role("alert"),
		)...,
	)
}
