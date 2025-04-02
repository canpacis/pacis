package components

import (
	icn "github.com/canpacis/pacis-ui/icons"
	r "github.com/canpacis/pacis-ui/renderer"
)

func Dialog(props ...r.I) r.Element {
	props = Join(
		props,
		X("data", "dialog"),
		X("trap.noscroll", "isOpen"),
		On("keydown.esc.window", "closeDialog"),
	)

	return r.Div(props...)
}

func DialogTrigger(trigger r.Element) r.Element {
	trigger.AddAttribute(On("click", "openDialog"))
	return trigger
}

func DialogContent(props ...r.I) r.Node {
	ps := []r.I{
		r.Class("fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 sm:rounded-lg"),
		r.Data(":state", "isOpen ? 'open' : 'closed'"),
		On("click.outside", "closeDialog"),
		X("show", "isOpen"),
		X("transition:enter.scale.96"),

		r.Div(
			r.Class("absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground"),

			Button(
				ButtonSizeIcon,
				ButtonVariantGhost,
				r.Class("h-6 w-6 rounded-sm"),
				On("click", "closeDialog"),

				icn.X(r.Class("h-4 w-4")),
			),
		),
	}
	ps = append(ps, props...)

	return r.Frag(
		// Overlay
		r.Div(
			r.Class("fixed inset-0 z-50 bg-black/80 w-dvw h-dvh data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"),
			r.Data(":state", "isOpen ? 'open' : 'closed'"),
			X("show", "isOpen"),
		),
		r.Div(ps...),
	)
}

func DialogHeader(props ...r.I) r.Element {
	props = Join(props, r.Class("flex flex-col space-y-1.5 text-center sm:text-left"))
	return r.Div(props...)
}

func DialogFooter(props ...r.I) r.Element {
	props = Join(props, r.Class("flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2"))
	return r.Div(props...)
}

func DialogTitle(props ...r.I) r.Element {
	props = Join(props, r.Class("text-lg font-semibold leading-none tracking-tight"))
	return r.Span(props...)
}

func DialogDescription(props ...r.I) r.Element {
	props = Join(props, r.Class("text-sm text-muted-foreground"))
	return r.Span(props...)
}
