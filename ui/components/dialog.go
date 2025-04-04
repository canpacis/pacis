package components

import (
	h "github.com/canpacis/pacis/ui/html"
	icn "github.com/canpacis/pacis/ui/icons"
)

func Dialog(props ...h.I) h.Element {
	props = Join(
		props,
		X("data", "dialog"),
		X("trap.noscroll", "isOpen"),
		On("keydown.esc.window", "closeDialog"),
	)

	return h.Div(props...)
}

func DialogTrigger(trigger h.Element) h.Element {
	trigger.AddAttribute(On("click", "openDialog"))
	return trigger
}

func DialogContent(props ...h.I) h.Node {
	ps := []h.I{
		h.Class("fixed left-[50%] top-[50%] z-50 grid w-full translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 rounded-lg"),
		h.Data(":state", "isOpen ? 'open' : 'closed'"),
		On("click.outside", "closeDialog"),
		X("show", "isOpen"),
		X("transition:enteh.scale.96"),

		h.Div(
			h.Class("absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground"),

			Button(
				ButtonSizeIcon,
				ButtonVariantGhost,
				h.Class("h-6 w-6 rounded-sm"),
				On("click", "closeDialog"),

				icn.X(h.Class("h-4 w-4")),
			),
		),
	}
	ps = append(ps, props...)

	return h.Frag(
		// Overlay
		h.Div(
			h.Class("fixed inset-0 z-50 bg-black/80 w-dvw h-dvh data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"),
			h.Data(":state", "isOpen ? 'open' : 'closed'"),
			X("show", "isOpen"),
		),
		h.Div(ps...),
	)
}

func DialogHeader(props ...h.I) h.Element {
	props = Join(props, h.Class("flex flex-col space-y-1.5 text-center sm:text-left"))
	return h.Div(props...)
}

func DialogFooter(props ...h.I) h.Element {
	props = Join(props, h.Class("flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2"))
	return h.Div(props...)
}

func DialogTitle(props ...h.I) h.Element {
	props = Join(props, h.Class("text-lg font-semibold leading-none tracking-tight"))
	return h.Span(props...)
}

func DialogDescription(props ...h.I) h.Element {
	props = Join(props, h.Class("text-sm text-muted-foreground"))
	return h.Span(props...)
}
