package separator

import (
	"log"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

type keytyp string

type Orientation int

func (o Orientation) Apply(el *html.Element) {
	el.Set(keytyp("group-orientation"), o)
}

func (Orientation) Done(el *html.Element) {
	o := el.Get(keytyp("group-orientation")).(Orientation)

	switch o {
	case Horizontal:
		el.ClassList.Add("h-[1px] w-full")
		el.SetAttribute("data-orientation", "horizontal")
	case Vertical:
		el.ClassList.Add("h-full w-[1px]")
		el.SetAttribute("data-orientation", "vertical")
	default:
		log.Fatalf("invalid group orientation: %d", o)
	}
}

func (Orientation) Item() {}

const (
	Horizontal = Orientation(iota)
	Vertical
)

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Role("none"),
			html.Data("slot", "separator"),
			html.Class("shrink-0 bg-border"),
			Horizontal,
		)...,
	)
}
