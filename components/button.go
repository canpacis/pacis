package components

import (
	"context"
	"io"

	r "github.com/canpacis/pacis/renderer"
)

type ButtonSize int

const (
	ButtonSizeDefault = ButtonSize(iota)
	ButtonSizeSm
	ButtonSizeLg
	ButtonSizeIcon
)

func (s ButtonSize) Render(ctx context.Context, w io.Writer) error {
	return r.Class(s.GetValue().(string)).Render(ctx, w)
}

func (ButtonSize) GetKey() string {
	return "class"
}

func (s ButtonSize) GetValue() any {
	switch s {
	case ButtonSizeDefault:
		return "h-9 px-4 py-2 has-[>svg]:px-3"
	case ButtonSizeSm:
		return "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5"
	case ButtonSizeLg:
		return "h-10 rounded-md px-6 has-[>svg]:px-4"
	case ButtonSizeIcon:
		return "size-9"
	default:
		panic("invalid button size property")
	}
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
	return r.Class(v.GetValue().(string)).Render(ctx, w)
}

func (ButtonVariant) GetKey() string {
	return "class"
}

func (v ButtonVariant) GetValue() any {
	switch v {
	case ButtonVariantDefault:
		return "shadow-xs bg-primary text-primary-foreground hover:bg-primary/90"
	case ButtonVariantDestructive:
		return "text-white shadow-xs bg-destructive hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60"
	case ButtonVariantOutline:
		return "border shadow-xs bg-background hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50"
	case ButtonVariantSecondary:
		return "shadow-xs bg-secondary text-secondary-foreground hover:bg-secondary/80"
	case ButtonVariantGhost:
		return "hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50"
	case ButtonVariantLink:
		return "text-primary underline-offset-4 hover:underline"
	default:
		panic("invalid button variant property")
	}
}

func Button(props ...r.I) r.Element {
	var variant ButtonVariant
	var size ButtonSize

	ps := []r.I{
		r.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"),
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
	return r.Btn(ps...)
}
