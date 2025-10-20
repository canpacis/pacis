package dropdown

import (
	"fmt"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/x"
)

func New(items ...html.Item) html.Node {
	return html.Span(
		components.ItemsOf(
			items,
			html.Data("slot", "dropdown-menu"),
			x.On("keydown", "open('keyboard', $refs.dropdownroot)", "down"),
			x.On("keydown", "open('keyboard', $refs.dropdownroot)", "up"),
			x.Data("dropdown"),
			x.Ref("dropdownroot"),
		)...,
	)
}

func Trigger(items ...html.Item) html.Node {
	return html.Button(components.ItemsOf(items, html.Data("slot", "dropdown-menu-trigger"), x.Ref("trigger"), Open)...)
}

func Group(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "dropdown-menu-group"), html.Role("group"))...)
}

func Label(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "dropdown-menu-label"), html.Class("px-2 py-1.5 text-sm font-semibold"))...)
}

func Separator(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "dropdown-menu-separator"), html.Role("separator"), html.Data("orientation", "horizontal"), html.Class("-mx-1 my-1 h-px bg-muted"))...)
}

func Shortcut(items ...html.Item) html.Node {
	return html.Div(components.ItemsOf(items, html.Data("slot", "dropdown-menu-shortcut"), html.Class("ml-auto text-xs tracking-widest opacity-60"))...)
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
	if el.Get("subtrigger") != nil {
		ref = "$refs.subtrigger"
		offset = 0
	}

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

		html.Div(
			components.ItemsOf(
				items,
				x.Cloak,
				x.X("trap.noscroll.noautofocus", "mouse"),
				x.X("trap.noscroll", "keyboard"),
				x.Show("opened"),
				x.Bind("data-state", "opened ? 'open' : 'closed'"),
				x.On("keydown", "$focus.next()", "down"),
				x.On("keydown", "$focus.previous()", "up"),
				x.On("keyup", "close($refs.dropdownroot)", "escape", x.Window),
				x.On("click", "close($refs.dropdownroot)", x.Outside),
				html.Data("slot", "dropdown-menu-content"),
				html.Class("z-50 min-w-[8rem] overflow-y-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"),
				Bottom,
				Offset(6),
				align,
			)...,
		),
	)
}

func Item(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			x.On("mouseenter", "$focus.focus($el)"),
			x.On("mouseleave", "$el.blur()"),
			x.On("keydown", "$el.click()", "enter"),
			x.On("keydown", "$el.click()", "space"),
			html.Role("menuitem"),
			html.Data("slot", "dropdown-menu-item"),
			html.Tabindex("0"),
			html.Class("focus:bg-accent focus:text-accent-foreground data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 dark:data-[variant=destructive]:focus:bg-destructive/20 data-[variant=destructive]:focus:text-destructive data-[variant=destructive]:*:[svg]:!text-destructive [&_svg:not([class*='text-'])]:text-muted-foreground relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[inset]:pl-8 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"),
			Close,
		)...,
	)
}

var Open = html.Attr("x-on:click", "open('mouse', $refs.dropdownroot)")
var Close = html.Attr("x-on:click", "close($refs.dropdownroot)")

func OpenOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "open('mouse', $refs.dropdownroot)")
}

func CloseOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "close($refs.dropdownroot)")
}
