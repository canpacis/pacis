package button

import (
	"log"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

type keytyp string

type Variant int

func (Variant) Item() {}

const (
	Default = Variant(iota)
	Destructive
	Outline
	Secondary
	Ghost
	Link
)

func (v Variant) Apply(el *html.Element) {
	el.Set(keytyp("variant"), v)
}

func (Variant) Done(el *html.Element) {
	v := el.Get(keytyp("variant")).(Variant)

	switch v {
	case Default:
		el.ClassList.Add("bg-primary text-primary-foreground shadow hover:bg-primary/90")
	case Destructive:
		el.ClassList.Add("bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/90")
	case Outline:
		el.ClassList.Add("border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground")
	case Secondary:
		el.ClassList.Add("bg-secondary text-secondary-foreground hover:bg-secondary/80")
	case Ghost:
		el.ClassList.Add("hover:bg-accent hover:text-accent-foreground")
	case Link:
		el.ClassList.Add("text-primary underline-offset-4 hover:underline")
	default:
		log.Fatalf("invalid button variant: %d", v)
	}
}

type Size int

func (Size) Item() {}

const (
	Md = Size(iota)
	Sm
	Lg
	Icon
)

func (s Size) Apply(el *html.Element) {
	el.Set(keytyp("size"), s)
}

func (Size) Done(el *html.Element) {
	s := el.Get(keytyp("size")).(Size)

	switch s {
	case Md:
		el.ClassList.Add("h-10 px-4 py-2")
	case Sm:
		el.ClassList.Add("h-9 rounded-md px-3")
	case Lg:
		el.ClassList.Add("h-11 rounded-md px-8")
	case Icon:
		el.ClassList.Add("h-10 w-10")
	default:
		log.Fatalf("invalid button size: %d", s)
	}
}

func New(items ...html.Item) html.Node {
	return html.Button(
		components.ItemsOf(
			items,
			html.Data("slot", "button"),
			html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0"),
			Default,
			Md,
		)...,
	)
}
