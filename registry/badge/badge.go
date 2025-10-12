package badge

import (
	"log"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

type Variant int

func (Variant) Item() {}

const (
	Default = Variant(iota)
	Destructive
	Outline
	Secondary
)

func (v Variant) Apply(el *html.Element) {
	el.Set("badge-variant", v)
}

func (Variant) Done(el *html.Element) {
	v := el.Get("badge-variant").(Variant)

	switch v {
	case Default:
		el.AddClass("border-transparent bg-primary text-primary-foreground hover:bg-primary/80")
	case Destructive:
		el.AddClass("border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80")
	case Outline:
		el.AddClass("text-foreground")
	case Secondary:
		el.AddClass("border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80")
	default:
		log.Fatalf("invalid badge variant: %d", v)
	}
}

func New(items ...html.Item) html.Node {
	return html.Span(
		components.ItemsOf(
			items,
			html.Data("slot", "badge"),
			html.Class("inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"),
			Default,
		)...,
	)
}
