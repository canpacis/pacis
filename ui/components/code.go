package components

import (
	"fmt"

	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func Code(code, lang string, props ...h.I) h.Element {
	return h.Div(
		Join(
			props,
			h.Class("relative"),

			Button(
				ButtonSizeIcon,
				ButtonVariantGhost,
				h.Class("!size-7 rounded-sm absolute top-2 right-2 md:top-3 md:right-3"),
				On("click", fmt.Sprintf("$clipboard(`%s`)", code)),

				icons.Clipboard(h.Class("size-3")),
				h.Span(h.Class("sr-only"), h.Text("Copy to Clipboard")),
			),
			h.Pre(
				h.Class("p-3 md:p-6 overflow-auto border rounded-lg"),

				h.Cde(h.Class("font-mono text-left break-normal leading-5 text-sm tabular-nums hyphens-none shadow-none font-medium whitespace-pre-wrap"), h.Class(lang), h.Text(code)),
			),
		)...,
	)
}
