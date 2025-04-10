package components

import (
	h "github.com/canpacis/pacis/ui/html"
)

var (
	BadgeVariantDefault = &GroupedClass{
		"badge-variant",
		"!border-transparent bg-primary text-primary-foreground [a&]:hover:bg-primary/90",
		true,
	}
	BadgeVariantSecondary = &GroupedClass{
		"badge-variant",
		"!border-transparent bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90",
		false,
	}
	BadgeVariantDestructive = &GroupedClass{
		"badge-variant",
		"!border-transparent bg-destructive text-white [a&]:hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60",
		false,
	}
	BadgeVariantOutline = &GroupedClass{
		"badge-variant",
		"text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
		false,
	}
)

func Badge(props ...h.I) h.Element {
	return h.Span(
		Join(
			props,
			BadgeVariantDefault,
			h.Class("inline-flex items-center justify-center rounded-md border px-2 py-0.5 text-xs font-medium w-fit whitespace-nowrap shrink-0 [&>svg]:size-3 gap-1 [&>svg]:pointer-events-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive transition-[color,box-shadow] overflow-hidden"),
		)...,
	)
}
