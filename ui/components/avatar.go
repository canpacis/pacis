package components

import (
	"fmt"

	h "github.com/canpacis/pacis/ui/html"
)

func AvatarFallback(props ...h.I) h.Element {
	props = Join(
		props,
		h.Class("bg-muted flex size-full items-center justify-center rounded-full text-sm absolute inset-0 z-10"),
	)
	return h.Div(props...)
}

func AvatarImage(src *h.Attribute, props ...h.I) h.Node {
	return h.Img(
		Join(
			props,
			X("data", fmt.Sprintf(`{url: "%s"}`, src.Value())),
			// D{"url": src.Value()},
			src,
			h.Class("aspect-square size-full relative z-20"),
			h.Attr(":class", "error ? 'hidden' : 'block'"),
			h.Attr(":src", "url"),
			On("error", "error = true"),
		)...,
	)
}

var (
	// AvatarSizeDefault = &GroupedClass{"avatar-size", "size-8", true}
	// AvatarSizeSm      = &GroupedClass{"avatar-size", "size-6 text-sm", false}
	// AvatarSizeLg      = &GroupedClass{"avatar-size", "size-12 text-lg", false}
	AvatarSizeDefault = h.Class("size-8")
	AvatarSizeSm      = h.Class("size-6 text-sm")
	AvatarSizeLg      = h.Class("size-12 text-lg")
)

func Avatar(props ...h.I) h.Element {
	return h.Div(Join(
		props,
		AvatarSizeDefault,
		h.Class("flex shrink-0 overflow-hidden rounded-full relative"),
		X("data", "{error: false}"),
		// D{"error": false},
	)...)
}
