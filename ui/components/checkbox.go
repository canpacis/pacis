package components

import (
	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func Checkbox(props ...h.I) h.Element {
	id := id()
	// TODO: route the checked property to input inside
	return h.Lbl(Join(props, h.Class("text-sm gap-2 items-center inline-flex w-fit-content"), h.HtmlFor(id),
		h.Span(
			h.Class("peer border-input dark:bg-input/30 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground dark:data-[state=checked]:bg-primary data-[state=checked]:border-primary focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive size-4 aspect-square shrink-0 rounded-[4px] border shadow-xs transition-shadow outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50"),
			D{"checked": false},
			h.Attr(":data-state", "checked ? 'checked' : 'unchecked'"),

			h.Span(
				h.Class("items-center justify-center text-current transition-none"),
				h.Attr(":class", "!checked ? 'hidden' : 'flex'"),
				icons.Check(h.Class("size-3.5")),
			),

			h.Inpt(
				h.ID(id),
				h.Type("checkbox"),
				h.Class("sr-only"),
				On("change", "checked = !checked, $dispatch('changed', { value: checked, event: $event })"),
			),
		),
	)...)
}
