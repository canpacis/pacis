package components

import (
	"io"

	r "github.com/canpacis/pacis/renderer"
)

type BadgeVariant int

const (
	BadgeVariantDefault = BadgeVariant(iota)
	BadgeVariantSecondary
	BadgeVariantDestructive
	BadgeVariantOutline
)

func (v BadgeVariant) Render(w io.Writer) error {
	class := ""

	switch v {
	case BadgeVariantDefault:
		class = "!border-transparent bg-primary text-primary-foreground [a&]:hover:bg-primary/90"
	case BadgeVariantSecondary:
		class = "!border-transparent bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90"
	case BadgeVariantDestructive:
		class = "!border-transparent bg-destructive text-white [a&]:hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60"
	case BadgeVariantOutline:
		class = "text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground"
	default:
		panic("invalid badge variant property")
	}

	return r.Class(class).Render(w)
}

func (BadgeVariant) Key() string {
	return "class"
}

func Badge(props ...r.Renderer) *r.Element {
	var variant BadgeVariant

	ps := []r.Renderer{
		r.Class("inline-flex items-center justify-center rounded-md border px-2 py-0.5 text-xs font-medium w-fit whitespace-nowrap shrink-0 [&>svg]:size-3 gap-1 [&>svg]:pointer-events-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive transition-[color,box-shadow] overflow-hidden"),
	}

	for _, prop := range props {
		switch prop := prop.(type) {
		case BadgeVariant:
			variant = prop
		default:
			ps = append(ps, prop)
		}
	}

	ps = append(ps, variant)
	return r.Span(ps...)
}
