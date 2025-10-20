package table

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

func New(items ...html.Item) html.Node {
	return html.Div(
		html.Class("relative w-full overflow-auto"),

		html.Table(
			components.ItemsOf(
				items,
				html.Class("w-full caption-bottom text-sm"),
			)...,
		),
	)
}

func Header(items ...html.Item) html.Node {
	return html.Thead(components.ItemsOf(items, html.Class("[&_tr]:border-b"))...)
}

func Body(items ...html.Item) html.Node {
	return html.Tbody(components.ItemsOf(items, html.Class("[&_tr:last-child]:border-0"))...)
}

func Footer(items ...html.Item) html.Node {
	return html.Tfoot(components.ItemsOf(items, html.Class("border-t bg-muted/50 font-medium [&>tr]:last:border-b-0"))...)
}

func Row(items ...html.Item) html.Node {
	return html.Tr(components.ItemsOf(items, html.Class("border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted"))...)
}

func Head(items ...html.Item) html.Node {
	return html.Th(components.ItemsOf(items, html.Class("h-12 px-4 text-left align-middle font-medium text-muted-foreground [&:has([role=checkbox])]:pr-0"))...)
}

func Cell(items ...html.Item) html.Node {
	return html.Td(components.ItemsOf(items, html.Class("p-4 align-middle [&:has([role=checkbox])]:pr-0"))...)
}

func Caption(items ...html.Item) html.Node {
	return html.Caption(components.ItemsOf(items, html.Class("mt-4 text-sm text-muted-foreground"))...)
}
