package components

import h "github.com/canpacis/pacis/ui/html"

func Table(props ...h.I) h.Element {
	return h.Div(
		h.Class("relative w-full overflow-auto"),

		h.Tble(
			Join(
				props,
				h.Class("w-full caption-bottom text-sm"),
			)...,
		),
	)
}

func TableHeader(props ...h.I) h.Element {
	return h.Th(Join(props, h.Class("[&_tr]:border-b p-2"))...)
}

func TableBody(props ...h.I) h.Element {
	return h.Tbody(Join(props, h.Class("[&_tr:last-child]:border-0"))...)
}

func TableRow(props ...h.I) h.Element {
	return h.Tr(
		Join(
			props,
			h.Class("border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted"),
		)...,
	)
}

func TableCell(props ...h.I) h.Element {
	return h.Td(
		Join(
			props,
			h.Class("p-2 align-middle [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]"),
		)...,
	)
}
