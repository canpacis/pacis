package accordion

import (
	"fmt"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/lucide"
	"github.com/canpacis/pacis/x"
)

func New(value string, items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			x.Data(fmt.Sprintf("accordion('%s')", value)),
			x.Ref("accordionroot"),
			html.Data("slot", "accordion"),
			html.Data("orientation", "vertical"),
			html.Class("h-fit w-full"),
		)...,
	)
}

func Item(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			x.Bind("data-state", "active === $el.getAttribute('data-value')"),
			html.Data("slot", "accordion-item"),
			html.Data("orientation", "vertical"),
			html.Class("border-b last:border-b-0"),
		)...,
	)
}

func Trigger(items ...html.Item) html.Node {
	items = components.ItemsOf(
		items,
		Select,
		html.Type("button"),
		x.Bind("data-state", "active === $el.parentElement.parentElement.getAttribute('data-value') ? 'open' : 'closed'"),
		x.Bind("aria-expanded", "active === $el.parentElement.parentElement.getAttribute('data-value')"),
		html.Data("slot", "accordion-trigger"),
		html.Data("orientation", "vertical"),
		html.Class("flex flex-1 items-center justify-between py-4 font-medium transition-all hover:underline [&[data-state=open]>svg]:rotate-180"),
	)
	items = append(items, lucide.ChevronDown(html.Class("h-4 w-4 shrink-0 transition-transform duration-200")))

	return html.H3(
		x.Bind("data-state", "active === $el.parentElement.getAttribute('data-value') ? 'open': 'closed'"),
		html.Class("flex"),

		html.Button(items...),
	)
}

func Content(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			x.Cloak,
			x.Show("active === $el.parentElement.getAttribute('data-value')"),
			x.Bind("data-state", "active === $el.parentElement.getAttribute('data-value') ? 'open': 'closed'"),
			html.Role("region"),
			html.Data("slot", "accordion-content"),
			html.Data("orientation", "vertical"),
			html.Class("overflow-hidden text-sm pb-4 pt-0 transition-alls"),
		)...,
	)
}

var Select = SelectOn("click")

func SelectOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "select($el.parentElement.parentElement.getAttribute('data-value'), $refs.accordionroot)")
}
