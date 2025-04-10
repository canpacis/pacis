package components

import (
	h "github.com/canpacis/pacis/ui/html"
)

var (
	ButtonSizeDefault = &GroupedClass{"button-size", "h-9 px-4 py-2 has-[>svg]:px-3", true}
	ButtonSizeSm      = &GroupedClass{"button-size", "h-8 rounded-md text-xs gap-1.5 px-3 has-[>svg]:px-2.5", false}
	ButtonSizeLg      = &GroupedClass{"button-size", "h-10 rounded-md px-6 has-[>svg]:px-4", false}
	ButtonSizeIcon    = &GroupedClass{"button-size", "size-9", false}
)

var (
	ButtonVariantDefault = &GroupedClass{
		"button-variant",
		"shadow-xs bg-primary text-primary-foreground hover:bg-primary/90",
		true,
	}
	ButtonVariantDestructive = &GroupedClass{
		"button-variant",
		"text-white shadow-xs bg-destructive hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40",
		false,
	}
	ButtonVariantOutline = &GroupedClass{
		"button-variant",
		"border shadow-xs bg-background hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50",
		false,
	}
	ButtonVariantSecondary = &GroupedClass{
		"button-variant",
		"shadow-xs bg-secondary text-secondary-foreground hover:bg-secondary/80",
		false,
	}
	ButtonVariantGhost = &GroupedClass{
		"button-variant",
		"hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50",
		false,
	}
	ButtonVariantLink = &GroupedClass{
		"button-variant",
		"text-primary underline-offset-4 hover:underline",
		false,
	}
)

func Button(props ...h.I) h.Element {
	props = Join(
		props,
		ButtonSizeDefault,
		ButtonVariantDefault,
		h.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"),
	)
	el := h.Btn(props...)

	attr, ok := el.GetAttribute("replace")
	if ok {
		el := attr.(*Replacer).element(props...)
		el.RemoveAttribute("replace")
		return el
	}
	return el
}
