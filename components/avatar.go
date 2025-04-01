package components

import (
	"bytes"
	"context"
	"fmt"
	"io"

	r "github.com/canpacis/pacis-ui/renderer"
)

func AvatarFallback(props ...r.I) r.Element {
	props = Join(
		props,
		r.Class("bg-muted flex size-full items-center justify-center rounded-full text-sm absolute inset-0 z-10"),
	)
	return r.Div(props...)
}

func AvatarImage(props ...r.I) r.Node {
	props = Join(
		props,
		r.Class("aspect-square size-full relative z-20"),
		r.Attr(":class", "error ? 'hidden' : 'block'"),
		r.Attr(":src", "url"),
		On("error", "error = true"),
	)

	el := r.Img(props...)

	url, ok := el.GetAttribute("src")
	if !ok {
		errset, ok := el.(r.ErrorSetter)
		if ok {
			errset.SetError(fmt.Errorf("avatar image component needs a src attribute"))
		} else {
			panic("avatar image component needs a src attribute")
		}
	} else {
		var buf bytes.Buffer
		url.Render(context.Background(), &buf)
		el.AddAttribute(D{"url": buf.String()})
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
	var value string

	switch v {
	case AvatarSizeDefault:
		value = "size-8"
	case AvatarSizeSm:
		value = "size-6 text-sm"
	case AvatarSizeLg:
		value = "size-12 text-lg"
	default:
		panic("invalid avatar size property")
	}

	return r.Class(value).Render(ctx, w)
}

func (AvatarSize) GetKey() string {
	return "class"
}

func (v AvatarSize) IsEmpty() bool {
	return false
}

func Avatar(props ...r.I) r.Element {
	var size AvatarSize

	ps := []r.I{
		r.Class("flex shrink-0 overflow-hidden rounded-full relative"),
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
