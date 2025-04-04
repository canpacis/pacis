package components

import (
	"context"
	"io"

	h "github.com/canpacis/pacis/ui/html"
)

type BadgeVariant int

const (
	BadgeVariantDefault = BadgeVariant(iota)
	BadgeVariantSecondary
	BadgeVariantDestructive
	BadgeVariantOutline
)

func (v BadgeVariant) Render(ctx context.Context, w io.Writer) error {
	var value string

	switch v {
	case BadgeVariantDefault:
		value = "!border-transparent bg-primary text-primary-foreground [a&]:hover:bg-primary/90"
	case BadgeVariantSecondary:
		value = "!border-transparent bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90"
	case BadgeVariantDestructive:
		value = "!border-transparent bg-destructive text-white [a&]:hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60"
	case BadgeVariantOutline:
		value = "text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground"
	default:
		panic("invalid badge variant property")
	}

	return h.Class(value).Render(ctx, w)
}

func (BadgeVariant) GetKey() string {
	return "class"
}

func (v BadgeVariant) IsEmpty() bool {
	return false
}

func Badge(props ...h.I) h.Element {
	var variant BadgeVariant

	ps := []h.I{
		h.Class("inline-flex items-center justify-center rounded-md border px-2 py-0.5 text-xs font-medium w-fit whitespace-nowrap shrink-0 [&>svg]:size-3 gap-1 [&>svg]:pointer-events-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive transition-[color,box-shadow] overflow-hidden"),
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
	return h.Span(ps...)
}
