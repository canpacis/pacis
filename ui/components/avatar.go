package components

import (
	"bytes"
	"context"
	"fmt"
	"io"

	h "github.com/canpacis/pacis/ui/html"
)

func AvatarFallback(props ...h.I) h.Element {
	props = Join(
		props,
		h.Class("bg-muted flex size-full items-center justify-center rounded-full text-sm absolute inset-0 z-10"),
	)
	return h.Div(props...)
}

func AvatarImage(props ...h.I) h.Node {
	props = Join(
		props,
		h.Class("aspect-square size-full relative z-20"),
		h.Attr(":class", "error ? 'hidden' : 'block'"),
		h.Attr(":src", "url"),
		On("error", "error = true"),
	)

	el := h.Img(props...)

	url, ok := el.GetAttribute("src")
	if !ok {
		errset, ok := el.(h.ErrorSetter)
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

	return h.Try(el, ErrorText)
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

	return h.Class(value).Render(ctx, w)
}

func (AvatarSize) GetKey() string {
	return "class"
}

func (v AvatarSize) IsEmpty() bool {
	return false
}

func Avatar(props ...h.I) h.Element {
	var size AvatarSize

	ps := []h.I{
		h.Class("flex shrink-0 overflow-hidden rounded-full relative"),
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

	return h.Div(ps...)
}
