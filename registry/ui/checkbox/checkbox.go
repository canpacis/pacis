package checkbox

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/lucide"
	"github.com/canpacis/pacis/x"
)

func New(items ...html.Item) html.Node {
	return html.Span(
		x.Data(map[string]any{"checked": false}),
		html.Class("size-4 inline-flex"),

		html.Button(
			components.ItemsOf(
				items,
				Toggle,

				x.Bind("aria-checked", "checked"),
				x.Bind("data-state", "checked ? 'checked' : 'unchecked'"),

				html.Value("on"),
				html.Data("slot", "checkbox"),
				html.Class("peer border-input dark:bg-input/30 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground dark:data-[state=checked]:bg-primary data-[state=checked]:border-primary focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive size-4 shrink-0 rounded-[4px] border shadow-xs transition-shadow outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50"),

				html.Span(
					x.Show("checked"),
					x.Bind("data-state", "checked ? 'checked' : 'unchecked'"),

					html.Data("slot", "checkbox-indicator"),
					html.Class("flex items-center justify-center text-current transition-none pointer-events-none"),

					lucide.Check(html.Class("h-4 w-4")),
				),
			)...,
		),
	)
}

func ToggleOn(event string) *html.Attribute {
	return x.On(event, "checked = !checked")
}

var Toggle = ToggleOn("click")
