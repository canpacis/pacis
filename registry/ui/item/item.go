package item

import (
	"fmt"

	"components/ui/separator"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

type Variant = components.Variant

const (
	DefaultVariant = Variant(iota)
	Outline
	Muted
)

var variant = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case DefaultVariant:
		el.AddClass("bg-transparent")
		el.SetAttribute("data-variant", "default")
	case Outline:
		el.AddClass("border-border")
		el.SetAttribute("data-variant", "outline")
	case Muted:
		el.AddClass("bg-muted/50")
		el.SetAttribute("data-variant", "muted")
	default:
		panic(fmt.Sprintf("invalid item variant: %d", v))
	}
})

type Size = components.Size

const (
	DefaultSize = Size(iota)
	Sm
)

var size = components.NewSizeApplier(func(el *html.Element, s components.Size) {
	switch s {
	case DefaultSize:
		el.AddClass("gap-4 p-4")
		el.SetAttribute("data-size", "default")
	case Sm:
		el.AddClass("gap-2.5 px-4 py-3")
		el.SetAttribute("data-size", "sm")
	default:
		panic(fmt.Sprintf("invalid item size: %d", s))
	}
})

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "item"),
			html.Class("group/item [a]:hover:bg-accent/50 focus-visible:border-ring focus-visible:ring-ring/50 [a]:transition-colors flex flex-wrap items-center rounded-md border border-transparent text-sm outline-none transition-colors duration-100 focus-visible:ring-[3px]"),
			DefaultVariant,
			DefaultSize,
			variant,
			size,
		)...,
	)
}

func Group(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Role("list"),
			html.Data("slot", "item-group"),
			html.Class("group/item-group flex flex-col"),
		)...,
	)
}

func Separator(items ...html.Item) html.Node {
	return separator.New(
		components.ItemsOf(
			items,
			html.Data("slot", "item-separator"),
			html.Class("my-0"),
			separator.Horizontal,
		)...,
	)
}

type MediaVariant = components.Variant

const (
	MediaDefault = MediaVariant(iota)
	MediaIcon
	MediaImage
)

var mediaVariant = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case MediaDefault:
		el.AddClass("bg-transparent")
		el.SetAttribute("data-variant", "default")
	case MediaIcon:
		el.AddClass("bg-muted size-8 rounded-sm border [&_svg:not([class*='size-'])]:size-4")
		el.SetAttribute("data-variant", "icon")
	case MediaImage:
		el.AddClass("size-10 overflow-hidden rounded-sm [&_img]:size-full [&_img]:object-cover")
		el.SetAttribute("data-variant", "image")
	default:
		panic(fmt.Sprintf("invalid item media variant: %d", v))
	}
})

func Media(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "item-media"),
			html.Class("flex shrink-0 items-center justify-center gap-2 group-has-[[data-slot=item-description]]/item:translate-y-0.5 group-has-[[data-slot=item-description]]/item:self-start [&_svg]:pointer-events-none"),
			MediaDefault,
			mediaVariant,
		)...)
}

func Content(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "item-content"),
			html.Class("flex flex-1 flex-col gap-1 [&+[data-slot=item-content]]:flex-none"),
		)...,
	)
}

func Title(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "item-title"),
			html.Class("flex w-fit items-center gap-2 text-sm font-medium leading-snug"),
		)...,
	)
}

func Description(items ...html.Item) html.Node {
	return html.P(
		components.ItemsOf(
			items,
			html.Data("slot", "item-description"),
			html.Class("text-muted-foreground line-clamp-2 text-balance text-sm font-normal leading-normal [&>a:hover]:text-primary [&>a]:underline [&>a]:underline-offset-4"),
		)...,
	)
}

func Actions(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "item-actions"),
			html.Class("flex items-center gap-2"),
		)...,
	)
}

func Header(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "item-header"),
			html.Class("flex basis-full items-center justify-between gap-2"),
		)...,
	)
}

func Footer(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "item-footer"),
			html.Class("flex basis-full items-center justify-between gap-2"),
		)...,
	)
}
