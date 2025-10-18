package empty

import (
	"fmt"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "empty"),
			html.Class("flex min-w-0 flex-1 flex-col items-center justify-center gap-6 text-balance rounded-lg border-dashed p-6 text-center md:p-12"),
		)...,
	)
}

func Header(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "empty-header"), html.Class("flex max-w-sm flex-col items-center gap-2 text-center"))...)
}

type MediaVariant = components.Variant

const (
	MediaDefault = MediaVariant(iota)
	MediaIcon
)

var mediaVariant = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case MediaDefault:
		el.AddClass("bg-transparent")
		el.SetAttribute("data-variant", "default")
	case MediaIcon:
		el.AddClass("bg-muted text-foreground flex size-10 shrink-0 items-center justify-center rounded-lg [&_svg:not([class*='size-'])]:size-6")
		el.SetAttribute("data-variant", "icon")
	default:
		panic(fmt.Sprintf("invalid empty media variant: %d", v))
	}
})

func Media(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "empty-icon"),
			html.Class("mb-2 flex shrink-0 items-center justify-center [&_svg]:pointer-events-none [&_svg]:shrink-0"),
			mediaVariant,
		)...,
	)
}

func Title(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "empty-title"), html.Class("text-lg font-medium tracking-tight"))...)
}

func Description(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "empty-description"), html.Class("text-muted-foreground [&>a:hover]:text-primary text-sm/relaxed [&>a]:underline [&>a]:underline-offset-4"))...)
}

func Content(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "empty-content"), html.Class("flex w-full min-w-0 max-w-sm flex-col items-center gap-4 text-balance text-sm"))...)
}
