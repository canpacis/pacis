package components

import h "github.com/canpacis/pacis/ui/html"

func Skeleton(props ...h.I) h.Element {
	return h.Div(
		Join(
			props,
			h.Class("animate-pulse rounded-md bg-primary/10"),
		)...,
	)
}
