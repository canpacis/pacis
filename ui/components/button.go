package components

import (
	"context"
	"io"

	h "github.com/canpacis/pacis/ui/html"
)

type ButtonSize int

const (
	ButtonSizeDefault = ButtonSize(iota)
	ButtonSizeSm
	ButtonSizeLg
	ButtonSizeIcon
)

func (s ButtonSize) Render(ctx context.Context, w io.Writer) error {
	var value string

	switch s {
	case ButtonSizeDefault:
		value = "h-9 px-4 py-2 has-[>svg]:px-3"
	case ButtonSizeSm:
		value = "h-8 rounded-md text-xs gap-1.5 px-3 has-[>svg]:px-2.5"
	case ButtonSizeLg:
		value = "h-10 rounded-md px-6 has-[>svg]:px-4"
	case ButtonSizeIcon:
		value = "size-9"
	default:
		panic("invalid button size property")
	}

	return h.Class(value).Render(ctx, w)
}

func (ButtonSize) GetKey() string {
	return "class"
}

func (s ButtonSize) IsEmpty() bool {
	return false
}

type ButtonVariant int

const (
	ButtonVariantDefault = ButtonVariant(iota)
	ButtonVariantDestructive
	ButtonVariantOutline
	ButtonVariantSecondary
	ButtonVariantGhost
	ButtonVariantLink
)

func (v ButtonVariant) Render(ctx context.Context, w io.Writer) error {
	var value string

	switch v {
	case ButtonVariantDefault:
		value = "shadow-xs bg-primary text-primary-foreground hover:bg-primary/90"
	case ButtonVariantDestructive:
		value = "text-white shadow-xs bg-destructive hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40"
	case ButtonVariantOutline:
		value = "border shadow-xs bg-background hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50"
	case ButtonVariantSecondary:
		value = "shadow-xs bg-secondary text-secondary-foreground hover:bg-secondary/80"
	case ButtonVariantGhost:
		value = "hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50"
	case ButtonVariantLink:
		value = "text-primary underline-offset-4 hover:underline"
	default:
		panic("invalid button variant property")
	}

	return h.Class(value).Render(ctx, w)
}

func (ButtonVariant) GetKey() string {
	return "class"
}

func (v ButtonVariant) IsEmpty() bool {
	return false
}

func Button(props ...h.I) h.Element {
	var variant ButtonVariant
	var size ButtonSize

	ps := []h.I{
		h.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"),
	}

	for _, prop := range props {
		switch prop := prop.(type) {
		case ButtonSize:
			size = prop
		case ButtonVariant:
			variant = prop
		default:
			ps = append(ps, prop)
		}
	}

	ps = append(ps, variant)
	ps = append(ps, size)

	el := h.Btn(ps...)

	attr, ok := el.GetAttribute("replace")
	if ok {
		el := attr.(*Replacer).element(ps...)
		el.RemoveAttribute("replace")
		return el
	}
	return el
}
