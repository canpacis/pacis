package components

import (
	"github.com/canpacis/pacis-ui/icons"
	r "github.com/canpacis/pacis-ui/renderer"
)

func Breadcrumb(props ...r.I) r.Element {
	return r.Div(
		r.Aria("label", "breadcrumb"),
		r.Ol(
			Join(
				props,
				r.Class("flex flex-wrap items-center gap-1.5 break-words text-sm text-muted-foreground sm:gap-2.5"),
			)...,
		),
	)
}

func BreadcrumbItem(props ...r.I) r.Element {
	return r.Li(
		Join(
			props,
			r.Class("inline-flex items-center gap-1.5"),
		)...,
	)
}

func BreadcrumbSeperator(props ...r.I) r.Element {
	return r.Li(
		Join(
			props,
			r.Role("presentation"),
			r.Aria("hidden", "true"),

			icons.ChevronRight(r.Class("size-3.5")),
		)...,
	)
}
