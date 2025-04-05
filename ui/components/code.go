package components

import (
	"fmt"

	. "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func Code(code, lang string, props ...I) Element {
	return Div(
		Join(
			props,
			Class("relative"),

			Button(
				ButtonSizeIcon,
				ButtonVariantGhost,
				Class("!size-7 rounded-sm absolute top-2 right-2 md:top-3 md:right-3"),
				On("click", fmt.Sprintf("$clipboard(`%s`)", code)),

				icons.Clipboard(Class("size-3")),
			),
			Pre(
				Class("p-3 md:p-6 overflow-auto border rounded-lg"),

				Cde(Class("font-mono text-left break-normal leading-5 text-sm tabular-nums hyphens-none shadow-none font-medium whitespace-pre-wrap"), Class(lang), Text(code)),
			),
		)...,
	)
}
