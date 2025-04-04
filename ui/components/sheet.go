package components

import (
	"context"
	"io"

	h "github.com/canpacis/pacis/ui/html"
)

type SheetVariant int

const (
	SheetVariantLeft = SheetVariant(iota)
	SheetVariantTop
	SheetVariantRight
	SheetVariantBottom
)

func (v SheetVariant) Render(ctx context.Context, w io.Writer) error {
	var value string

	switch v {
	case SheetVariantLeft:
		value = "top-0 bottom-0 left-0 h-dvh w-3/4 border-r data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left sm:max-w-sm"
	case SheetVariantTop:
		value = "left-0 right-0 top-0 border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top"
	case SheetVariantRight:
		value = "top-0 bottom-0 right-0 h-dvh w-3/4 border-l data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right sm:max-w-sm"
	case SheetVariantBottom:
		value = "left-0 right-0 bottom-0 border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom"
	default:
		panic("invalid sheet variant property")
	}

	return h.Class(value).Render(ctx, w)
}

func (SheetVariant) GetKey() string {
	return "class"
}

func (v SheetVariant) IsEmpty() bool {
	return false
}

func Sheet(props ...h.I) h.Element {
	return h.Div(Join(props, X("data", "sheet"))...)
}

func SheetTrigger(trigger h.Element) h.Element {
	trigger.AddAttribute(On("click", "openSheet"))
	return trigger
}

func SheetContent(props ...h.I) h.Node {
	var variant SheetVariant
	contentprops := []h.I{
		X("show", "isOpen"),
		X("trap.noscroll", "isOpen"),
		On("click.outside", "closeSheet"),
		h.Data(":state", "isOpen ? 'opened' : 'closed'"),
		h.Class("fixed z-50 gap-4 bg-background p-6 shadow-lg"),
	}

	for _, prop := range props {
		switch prop := prop.(type) {
		case SheetVariant:
			variant = prop
		default:
			contentprops = append(contentprops, prop)
		}
	}
	contentprops = append(contentprops, variant)

	return h.Frag(
		// Overlay
		h.Div(
			X("cloak"),
			X("show", "isOpen"),
			h.Data(":state", "isOpen ? 'opened' : 'closed'"),
			h.Class("fixed h-dvh inset-0 z-50 bg-black/80"),
			X("transison:enter", "transition duration-300"),
			X("transison:enter-start", "fade-out-0"),
			X("transison:enter-end", "fade-out-1"),
			X("transison:leave", "transition duration-300"),
			X("transison:leave-start", "fade-out-0"),
			X("transison:leave-end", "fade-out-1"),
		),
		h.Div(contentprops...),
	)
}
