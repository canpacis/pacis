package switcher

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/x"
)

func New(items ...html.Item) html.Node {
	return html.Button(
		components.ItemsOf(
			items,

			x.Data(map[string]any{"checked": false}),
			x.Bind("aria-checked", "checked"),
			x.Bind("data-state", "checked ? 'checked' : 'unchecked'"),

			Toggle,

			html.Type("button"),
			html.Value("on"),
			html.Data("slot", "switch"),
			html.Role("switch"),
			html.Class("peer inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-primary data-[state=unchecked]:bg-input"),

			html.Span(
				x.Bind("data-state", "checked ? 'checked' : 'unchecked'"),
				html.Data("slot", "switch-thumb"),
				html.Class("pointer-events-none block h-5 w-5 rounded-full bg-background shadow-lg ring-0 transition-transform data-[state=checked]:translate-x-5 data-[state=unchecked]:translate-x-0"),
			),
		)...,
	)
}

func ToggleOn(event string) *html.Attribute {
	return x.On(event, "checked = !checked")
}

var Toggle = ToggleOn("click")
