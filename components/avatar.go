package components

import (
	"context"
	"io"

	r "github.com/canpacis/pacis/renderer"
)

func AvatarFallback(props ...r.I) r.Element {
	ps := []r.I{
		r.Class("bg-muted flex size-full items-center justify-center rounded-full"),
		r.Attr(":class", "!error ? 'hidden' : 'block'"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func AvatarImage(props ...r.I) r.Element {
	ps := []r.I{
		r.Class("aspect-square size-full"),
		r.Attr(":class", "error ? 'hidden' : 'block'"),
		r.Attr(":src", "url"),
		On("error", "error = true"),
	}
	ps = append(ps, props...)

	el := r.Img(ps...)
	src, ok := el.GetAttribute("src")
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

func (v AvatarSize) Render(ctx context.Context, w io.Writer) error {
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

	return r.Class(class).Render(ctx, w)
}

func (AvatarSize) Key() string {
	return "class"
}

func Avatar(props ...r.I) r.Element {
	var size AvatarSize

	ps := []r.I{
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
