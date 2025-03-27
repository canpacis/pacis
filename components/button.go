package components

import (
	"io"

	"github.com/canpacis/pacis/renderer"
)

type ButtonSize int

const (
	ButtonSizeDefault = ButtonSize(iota)
	ButtonSizeSm
	ButtonSizeLg
	ButtonSizeIcon
)

func (s ButtonSize) Render(w io.Writer) error {
	class := ""

	switch s {
	case ButtonSizeDefault:
		class = "h-9 px-4 py-2 has-[>svg]:px-3"
	case ButtonSizeSm:
		class = "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5"
	case ButtonSizeLg:
		class = "h-10 rounded-md px-6 has-[>svg]:px-4"
	case ButtonSizeIcon:
		class = "size-9"
	default:
		panic("invalid button size property")
	}

	return renderer.Class(class).Render(w)
}

func (ButtonSize) Key() string {
	return "class"
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

func (v ButtonVariant) Render(w io.Writer) error {
	class := ""

	switch v {
	case ButtonVariantDefault:
		class = "bg-primary text-primary-foreground shadow-xs hover:bg-primary/90"
	case ButtonVariantDestructive:
		class = "bg-destructive text-white shadow-xs hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60"
	case ButtonVariantOutline:
		class = "border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50"
	case ButtonVariantSecondary:
		class = "bg-secondary text-secondary-foreground shadow-xs hover:bg-secondary/80"
	case ButtonVariantGhost:
		class = "hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50"
	case ButtonVariantLink:
		class = "text-primary underline-offset-4 hover:underline"
	default:
		panic("invalid button variant property")
	}

	return renderer.Class(class).Render(w)
}

func (ButtonVariant) Key() string {
	return "class"
}

func Button(props ...renderer.Renderer) *renderer.Element {
	var variant ButtonVariant
	var size ButtonSize

	ps := []renderer.Renderer{
		renderer.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"),
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
	return renderer.Btn(ps...)
}
