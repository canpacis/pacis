package alert

import (
	"log"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

type Variant = components.Variant

const (
	Default = Variant(iota)
	Destructive
)

var variant = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case Default:
		el.AddClass("bg-background text-foreground")
	case Destructive:
		el.AddClass("border-destructive/50 text-destructive dark:border-destructive [&>svg]:text-destructive")
	default:
		log.Fatalf("invalid alert variant: %d", v)
	}
})

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Role("alert"),
			html.Data("slot", "alert"),
			html.Class("relative w-full rounded-lg border p-4 [&>svg~*]:pl-7 [&>svg+div]:translate-y-[-3px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg]:text-foreground"),
			Default,
			variant,
		)...,
	)
}

func Title(items ...html.Item) html.Node {
	return html.H5(
		components.ItemsOf(
			items,
			html.Data("slot", "alert-title"),
			html.Class("mb-1 font-medium leading-none tracking-tight"),
		)...,
	)
}

func Description(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "alert-description"),
			html.Class("text-sm [&_p]:leading-relaxed"),
		)...,
	)
}
