package components

import (
	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

/*
	A window overlaid on either the primary window or another dialog window, rendering the content underneath inert.

Usage:

	Dialog(
		DialogTrigger(
			Button(Text("Open Dialog")),
		),
		DialogContent(
			DialogHeader(
				DialogTitle(Text("Are you absolutely sure?")),
				DialogDescription(Text("This action cannot be undone. This will permanently delete your account and remove your data from our servers.")),
			),
		),
	)
*/
func Dialog(props ...h.I) h.Element {
	props = Join(
		props,
		X("data", "dialog"),
		X("trap.noscroll", "open"),
		DismissDialogOn("keydown.esc.window"),
	)

	return h.Div(props...)
}

// The trigger slot for the dialog component
func DialogTrigger(trigger h.Element) h.Element {
	trigger.AddAttribute(OpenDialog)
	return trigger
}

// The content slot for the dialog component
func DialogContent(props ...h.I) h.Node {
	ps := []h.I{
		h.Class("fixed left-[50%] top-[50%] z-50 grid w-full translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 rounded-lg"),
		h.Data(":state", "open ? 'open' : 'closed'"),
		DismissDialogOn("click.outside"),
		X("show", "open"),
		X("transition:enteh.scale.96"),

		h.Div(
			h.Class("absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground"),

			Button(
				ButtonSizeIcon,
				ButtonVariantGhost,
				DismissDialog,
				h.Class("h-6 w-6 rounded-sm"),

				icons.X(h.Class("h-4 w-4")),
			),
		),
	}
	ps = append(ps, props...)

	return h.Frag(
		// Overlay
		h.Div(
			h.Class("fixed inset-0 z-50 bg-black/80 w-dvw h-dvh data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"),
			h.Data(":state", "open ? 'open' : 'closed'"),
			X("show", "open"),
		),
		h.Div(ps...),
	)
}

// The header slot for the dialog component
func DialogHeader(props ...h.I) h.Element {
	props = Join(props, h.Class("flex flex-col space-y-1.5 text-center sm:text-left"))
	return h.Div(props...)
}

// The footer slot for the dialog component
func DialogFooter(props ...h.I) h.Element {
	props = Join(props, h.Class("flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2"))
	return h.Div(props...)
}

// The title slot for the dialog component
func DialogTitle(props ...h.I) h.Element {
	props = Join(props, h.Class("text-lg font-semibold leading-none tracking-tight"))
	return h.Span(props...)
}

// The description slot for the dialog component
func DialogDescription(props ...h.I) h.Element {
	props = Join(props, h.Class("text-sm text-muted-foreground"))
	return h.Span(props...)
}

func OpenDialogOn(event string) h.Attribute {
	return On(event, "openDialog()")
}

var OpenDialog = OpenDialogOn("click")

func CloseDialogOn(event, value string) h.Attribute {
	return On(event, fn("closeDialog", value))
}

func CloseDialog(value string) h.Attribute {
	return CloseDialogOn("click", value)
}

func DismissDialogOn(event string) h.Attribute {
	return On(event, "dismissDialog()")
}

var DismissDialog = DismissDialogOn("click")
