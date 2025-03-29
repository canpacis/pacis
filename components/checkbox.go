package components

import (
	icn "github.com/canpacis/pacis-ui/icons"
	r "github.com/canpacis/pacis-ui/renderer"
)

func Checkbox(label ...r.Text) r.Element {
	ps := []r.I{
		r.Class("text-sm gap-2 items-center inline-flex w-fit-content"),
		r.Span(
			r.Class("peer border-input dark:bg-input/30 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground dark:data-[state=checked]:bg-primary data-[state=checked]:border-primary focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive size-4 aspect-square shrink-0 rounded-[4px] border shadow-xs transition-shadow outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50"),
			D{"checked": false},
			r.Attr(":data-state", "checked ? 'checked' : 'unchecked'"),

			r.Span(
				r.Class("items-center justify-center text-current transition-none"),
				r.Attr(":class", "!checked ? 'hidden' : 'flex'"),
				icn.Check(r.Class("size-3.5")),
			),

			r.Inpt(
				r.Type("checkbox"),
				r.Class("sr-only"),
				On("change", "checked = !checked;"),
			),
		),
	}
	for _, label := range label {
		ps = append(ps, label)
	}

	return r.Lbl(ps...)
}
