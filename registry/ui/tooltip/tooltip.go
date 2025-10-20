package tooltip

import (
	"fmt"
	"time"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/x"
)

func New(delay time.Duration, items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			x.Ref("tooltiproot"),
			x.Data(fmt.Sprintf("tooltip(%d)", delay.Milliseconds())),

			html.Class("size-fit"),
		)...,
	)
}

func Trigger(items ...html.Item) html.Node {
	return html.Button(
		components.ItemsOf(
			items,

			x.Bind("data-state", "opened ? 'open' : 'close'"),
			x.Ref("trigger"),

			OpenOn("mouseenter"),
			CloseOn("mouseleave"),
		)...,
	)
}

type Align = components.Variant

const (
	Bottom = Align(iota)
	BottomStart
	BottomEnd
	Top
	TopStart
	TopEnd
	Left
	LeftStart
	LeftEnd
	Right
	RightStart
	RightEnd
)

type Offset = components.Size

var align = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	var ref = "$refs.trigger"
	var offset = el.Get("size").(components.Size)

	switch v {
	case Bottom:
		el.SetAttribute(fmt.Sprintf("x-anchor.bottom.offset.%d", offset), ref)
	case BottomStart:
		el.SetAttribute(fmt.Sprintf("x-anchor.bottom-start.offset.%d", offset), ref)
	case BottomEnd:
		el.SetAttribute(fmt.Sprintf("x-anchor.bottom-end.offset.%d", offset), ref)
	case Top:
		el.SetAttribute(fmt.Sprintf("x-anchor.top.offset.%d", offset), ref)
	case TopStart:
		el.SetAttribute(fmt.Sprintf("x-anchor.top-start.offset.%d", offset), ref)
	case TopEnd:
		el.SetAttribute(fmt.Sprintf("x-anchor.top-end.offset.%d", offset), ref)
	case Left:
		el.SetAttribute(fmt.Sprintf("x-anchor.left.offset.%d", offset), ref)
	case LeftStart:
		el.SetAttribute(fmt.Sprintf("x-anchor.left-start.offset.%d", offset), ref)
	case LeftEnd:
		el.SetAttribute(fmt.Sprintf("x-anchor.left-end.offset.%d", offset), ref)
	case Right:
		el.SetAttribute(fmt.Sprintf("x-anchor.right.offset.%d", offset), ref)
	case RightStart:
		el.SetAttribute(fmt.Sprintf("x-anchor.right-start.offset.%d", offset), ref)
	case RightEnd:
		el.SetAttribute(fmt.Sprintf("x-anchor.right-end.offset.%d", offset), ref)
	default:
		panic(fmt.Sprintf("invalid align variant: %d", v))
	}
})

func Content(items ...html.Item) html.Node {
	return html.Template(
		x.Teleport("body"),

		html.Span(
			components.ItemsOf(
				items,
				html.Attr("x-transition:enter", "transition-opacity"),
				html.Attr("x-transition:enter-start", "opacity-0"),
				html.Attr("x-transition:enter-end", "opacity-100"),
				html.Attr("x-transition:leave", "transition-opacity"),
				html.Attr("x-transition:leave-start", "opacity-100"),
				html.Attr("x-transition:leave-end", "opacity-0"),
				html.Class("z-50 overflow-hidden rounded-md border bg-popover px-3 py-1.5 text-sm text-popover-foreground shadow-md"),
				x.Show("opened"),
				Top,
				Offset(6),
				align,
			)...,
		),
	)
}

func OpenOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "open($refs.tooltiproot)")
}

func CloseOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "close($refs.tooltiproot)")
}
