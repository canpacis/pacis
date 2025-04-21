package components

import (
	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func ToastContainer(props ...h.I) h.I {
	return h.Div(
		Join(
			props,
			h.Class("fixed top-0 left-0 p-6 w-dvw h-dvh flex justify-end items-end size-fit z-50 pointer-events-none"),

			h.Div(
				h.Class("flex flex-col gap-2"),

				h.Template(
					X("for", "toast in $store.toast.visibleToasts"),
					h.Div(
						X("data", "toast"),
						X("show", "show"),
						X("bind:key", "toast.id"),
						X("transition.delay"),
						h.Class("pointer-events-auto relative border rounded-md p-4 min-w-90 bg-background text-sm transition-opacity"),

						h.P(h.Class("mb-1"), Textx("toast.content.title")),
						h.Span(h.Class("text-muted-foreground"), Textx("toast.content.message")),
						Button(
							On("click", "$store.toast.clear(toast.id)"),
							h.Class("absolute top-2 right-2 w-6 h-6"),
							ButtonSizeIcon,
							ButtonVariantGhost,

							icons.X(h.Class("size-3")),
						),
					),
				),
			),
		)...,
	)
}
