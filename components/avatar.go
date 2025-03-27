package components

import (
	"io"

	"github.com/canpacis/pacis/renderer"
)

func AvatarFallback(props ...renderer.Renderer) *renderer.Element {
	ps := []renderer.Renderer{
		renderer.Class("bg-muted flex size-full items-center justify-center rounded-full"),
		renderer.Attr(":class", "!error ? 'hidden' : 'block'"),
	}
	ps = append(ps, props...)

	return renderer.Div(ps...)
}

func AvatarImage(props ...renderer.Renderer) *renderer.Element {
	ps := []renderer.Renderer{
		renderer.Class("aspect-square size-full"),
		renderer.Attr(":class", "error ? 'hidden' : 'block'"),
		renderer.Attr(":src", "url"),
		renderer.Attr("@error", "error = true"),
	}
	ps = append(ps, props...)

	el := renderer.Img(ps...)
	src, ok := renderer.GetAttr(el, "src")
	if !ok {
		panic("avatar image component needs a src attribute")
	}

	ps = append(ps, D{"url": src})
	el = renderer.Img(ps...)
	return el
}

type AvatarSize int

const (
	AvatarSizeDefault = AvatarSize(iota)
	AvatarSizeSm
	AvatarSizeLg
)

func (v AvatarSize) Render(w io.Writer) error {
	class := ""

	switch v {
	case AvatarSizeDefault:
		class = "size-8"
	case AvatarSizeSm:
		class = "size-6"
	case AvatarSizeLg:
		class = "size-12"
	default:
		panic("invalid avatar size property")
	}

	return renderer.Class(class).Render(w)
}

func (AvatarSize) Key() string {
	return "class"
}

func Avatar(props ...renderer.Renderer) *renderer.Element {
	var size AvatarSize

	ps := []renderer.Renderer{
		renderer.Class("relative flex shrink-0 overflow-hidden rounded-full"),
		renderer.Attr("x-data", "{ 'error': false }"),
	}

	for _, prop := range props {
		switch prop := prop.(type) {
		case AvatarSize:
			size = prop
		default:
			ps = append(ps, prop)
		}
	}
	ps = append(ps, size)

	return renderer.Div(ps...)
}
