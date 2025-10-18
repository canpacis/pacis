package badge

import (
	"fmt"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

type Variant = components.Variant

const (
	Default = Variant(iota)
	Destructive
	Outline
	Secondary
)

var variant = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case Default:
		el.AddClass("border-transparent bg-primary text-primary-foreground hover:bg-primary/80")
	case Destructive:
		el.AddClass("border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80")
		fmt.Println(el)
	case Outline:
		el.AddClass("text-foreground")
	case Secondary:
		el.AddClass("border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80")
	default:
		panic(fmt.Sprintf("invalid badge variant: %d", v))
	}
})

func New(items ...html.Item) html.Node {
	return html.Span(
		components.ItemsOf(
			items,
			html.Data("slot", "badge"),
			html.Class("inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"),
			Default,
			variant,
		)...,
	)
}
