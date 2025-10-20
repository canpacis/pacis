package button

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
	Ghost
	Link
)

var variant = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case Default:
		el.AddClass("bg-primary text-primary-foreground hover:bg-primary/90")
	case Destructive:
		el.AddClass("bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/90")
	case Outline:
		el.AddClass("border border-input bg-background hover:bg-accent hover:text-accent-foreground")
	case Secondary:
		el.AddClass("bg-secondary text-secondary-foreground hover:bg-secondary/80")
	case Ghost:
		el.AddClass("hover:bg-accent hover:text-accent-foreground")
	case Link:
		el.AddClass("text-primary underline-offset-4 hover:underline")
	default:
		panic(fmt.Sprintf("invalid button variant: %d", v))
	}
})

type Size = components.Size

const (
	Md = Size(iota)
	Sm
	Lg
	Icon
)

var size = components.NewSizeApplier(func(el *html.Element, v components.Size) {
	switch v {
	case Md:
		el.AddClass("h-10 px-4 py-2")
	case Sm:
		el.AddClass("h-9 rounded-md px-3")
	case Lg:
		el.AddClass("h-11 rounded-md px-8")
	case Icon:
		el.AddClass("h-10 w-10")
	default:
		panic(fmt.Sprintf("invalid button size: %d", v))
	}
})

func New(items ...html.Item) html.Node {
	return html.Button(
		components.ItemsOf(
			items,
			html.Data("slot", "button"),
			html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0"),
			Default,
			Md,
			variant,
			size,
		)...,
	)
}
