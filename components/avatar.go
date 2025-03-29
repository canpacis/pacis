package components

import (
	"context"
	"fmt"
	"io"

	r "github.com/canpacis/pacis-ui/renderer"
)

func AvatarFallback(props ...r.I) r.Element {
	props = join(
		props,
		r.Class("bg-muted flex size-full items-center justify-center rounded-full text-sm"),
		r.Attr(":class", "!error ? 'hidden' : 'block'"),
	)
	return r.Div(props...)
}

func AvatarImage(props ...r.I) r.Node {
	props = join(
		props,
		r.Class("aspect-square size-full"),
		r.Attr(":class", "error ? 'hidden' : 'block'"),
		r.Attr(":src", "url"),
		On("error", "error = true"),
	)

	el := r.Img(props...)

	src, ok := el.GetAttribute("src")
	if !ok {
		errset, ok := el.(r.ErrorSetter)
		if ok {
			errset.SetError(fmt.Errorf("avatar image component needs a src attribute"))
		} else {
			panic("avatar image component needs a src attribute")
		}
	} else {
		el.AddAttribute(D{"url": src.GetValue()})
	}

	return r.Try(el, ErrorText)
}

type AvatarSize int

const (
	AvatarSizeDefault = AvatarSize(iota)
	AvatarSizeSm
	AvatarSizeLg
)

func (v AvatarSize) Render(ctx context.Context, w io.Writer) error {
	return r.Class(v.GetValue().(string)).Render(ctx, w)
}

func (AvatarSize) GetKey() string {
	return "class"
}

func (v AvatarSize) GetValue() any {
	switch v {
	case AvatarSizeDefault:
		return "size-8"
	case AvatarSizeSm:
		return "size-6"
	case AvatarSizeLg:
		return "size-12"
	default:
		panic("invalid avatar size property")
	}
}

func Avatar(props ...r.I) r.Element {
	var size AvatarSize

	ps := []r.I{
		r.Class("flex shrink-0 overflow-hidden rounded-full"),
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
