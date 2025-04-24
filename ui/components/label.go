package components

import h "github.com/canpacis/pacis/ui/html"

func Label(text string, props ...h.I) h.Element {
	return h.Lbl(
		Join(props,
			h.Class("text-sm w-full font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 block flex-1 [&>div>input]:mt-2"),
			h.Span(h.Class("inline-block"), h.Text(text)),
		)...,
	)
}
