package tabs

import (
	"fmt"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/x"
)

func New(value string, items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			x.Data(fmt.Sprintf("tabs('%s')", value)),
			x.Ref("tabsroot"),
			html.Data("slot", "tabs"),
			html.Aria("orientation", "horizontal"),
			html.Class("flex flex-col gap-2"),
		)...,
	)
}

func List(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Role("tablist"),
			html.Data("slot", "tabs-list"),
			html.Aria("orientation", "horizontal"),
			html.Class("bg-muted text-muted-foreground inline-flex h-9 w-fit items-center justify-center rounded-lg p-[3px]"),
		)...,
	)
}

func Trigger(items ...html.Item) html.Node {
	return html.Button(
		components.ItemsOf(
			items,
			x.Bind("aria-selected", "active === $el.getAttribute('data-value')"),
			x.Bind("data-state", "active === $el.getAttribute('data-value') ? 'active' : 'inactive'"),
			x.On("click", "select($el.getAttribute('data-value'), $refs.tabsroot)"),
			html.Class("inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-sm"),
			html.Type("button"),
			html.Role("tab"),
		)...,
	)
}

func Content(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			x.Show("active === $el.getAttribute('data-value')"),
			x.Bind("data-state", "active === $el.getAttribute('data-value') ? 'active' : 'inactive'"),
			x.Cloak,
			html.Role("tabpanel"),
			html.Data("slot", "tabs-content"),
			html.Aria("orientation", "horizontal"),
			html.Class("mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"),
		)...,
	)
}

func Select(value string) *html.Attribute {
	return SelectOn("click", value)
}

func SelectOn(event, value string) *html.Attribute {
	return x.On(event, fmt.Sprintf("select('%s', $refs.tabsroot)", value))
}
