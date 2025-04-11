package components

import (
	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func Textarea(props ...h.I) h.Element {
	return h.Lbl(
		h.Class("relative"),

		h.Txtarea(
			Join(
				props,
				h.Class("flex min-h-[60px] w-full rounded-md border border-input bg-background px-3 py-2 text-base shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"),
			)...,
		),
		h.Span(
			h.Class("absolute bottom-px right-px pr-0.5 rounded-sm pointer-events-none bg-background size-fit"),
			icons.GripHorizontal(h.Class("size-4 stroke-1 text-muted-foreground")),
		),
	)
}
