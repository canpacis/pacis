package components

import (
	"fmt"

	. "github.com/canpacis/pacis/ui/components"
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
				Class("size-7 rounded-sm absolute top-2 right-2 md:top-4 md:right-4"),
				On("click", fmt.Sprintf("$clipboard(`%s`)", code)),

				icons.Clipboard(Class("size-3")),
			),
			Pre(
				Cde(Class(lang), Text(code)),
			),
		)...,
	)
}
