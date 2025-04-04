package components

import (
	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func Breadcrumb(props ...h.I) h.Element {
	return h.Div(
		h.Aria("label", "breadcrumb"),
		h.Ol(
			Join(
				props,
				h.Class("flex flex-wrap items-center gap-1.5 break-words text-sm text-muted-foreground sm:gap-2.5"),
			)...,
		),
	)
}

func BreadcrumbItem(props ...h.I) h.Element {
	return h.Li(
		Join(
			props,
			h.Class("inline-flex items-center gap-1.5"),
		)...,
	)
}

func BreadcrumbSeperator(props ...h.I) h.Element {
	return h.Li(
		Join(
			props,
			h.Role("presentation"),
			h.Aria("hidden", "true"),

			icons.ChevronRight(h.Class("size-3.5")),
		)...,
	)
}
