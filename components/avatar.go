package components

import (
	"io"

	r "github.com/canpacis/pacis/renderer"
)

func AvatarFallback(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("bg-muted flex size-full items-center justify-center rounded-full"),
		r.Attr(":class", "!error ? 'hidden' : 'block'"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func AvatarImage(props ...r.Renderer) *r.Element {
	ps := []r.Renderer{
		r.Class("aspect-square size-full"),
		r.Attr(":class", "error ? 'hidden' : 'block'"),
		r.Attr(":src", "url"),
		On("error", "error = true"),
	}
	ps = append(ps, props...)

	el := r.Img(ps...)
	src, ok := r.GetAttr(el, "src")
	if !ok {
		panic("avatar image component needs a src attribute")
	}

	ps = append(ps, D{"url": src})
	el = r.Img(ps...)
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

	return r.Class(class).Render(w)
}

func (AvatarSize) Key() string {
	return "class"
}

func Avatar(props ...r.Renderer) *r.Element {
	var size AvatarSize

	ps := []r.Renderer{
		r.Class("relative flex shrink-0 overflow-hidden rounded-full"),
		D{"error": false},
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

	return r.Div(ps...)
}
