package selector

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
			x.Data(fmt.Sprintf("select('%s')", value)),
			x.On("keydown", "open('keyboard', $refs.selectroot)", "down"),
			x.On("keydown", "open('keyboard', $refs.selectroot)", "up"),
			x.Ref("selectroot"),
			html.Class("relative"),
		)...,
	)
}

func Trigger(items ...html.Item) html.Node {
	items = components.ItemsOf(
		items,
		Open,
		x.Bind("aria-expanded", "opened"),
		x.Bind("data-state", "opened ? 'open' : 'closed'"),
		x.Ref("trigger"),
		html.Type("button"),
		html.Role("combobox"),
		html.Autocomplete("none"),
		html.Data("slot", "select-trigger"),
		html.Class("flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background data-[placeholder]:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1"),
	)
	items = append(items, html.Aria("hidden", "true"), lucide.ChevronDown(html.Class("h-4 w-4 opacity-50")))

	return html.Button(items...)
}

func Value(placeholder string, items ...html.Item) html.Node {
	return html.Span(
		components.ItemsOf(
			items,
			html.Data("slot", "select-value"),
			html.Class("pointer-events-none"),

			x.Text(fmt.Sprintf("label($refs.selectroot, '%s')", placeholder)),
		)...,
	)
}

func Group(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "select-group"),
			html.Role("group"),
		)...,
	)
}

func Label(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "select-label"),
			html.Class("text-muted-foreground px-2 py-1.5 text-xs"),
		)...,
	)
}

func Separator(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Aria("hidden", "true"),
			html.Data("slot", "select-separator"),
			html.Class("-mx-1 my-1 h-px bg-muted"),
		)...,
	)
}

func Content(items ...html.Item) html.Node {
	return html.Div(
		x.Cloak,
		x.Show("opened"),
		x.X("trap.noautofocus", "mouse"),
		x.X("trap", "keyboard"),
		x.On("keydown", "$focus.next()", "down"),
		x.On("keydown", "$focus.previous()", "up"),
		x.On("keyup", "close($refs.selectroot)", "escape", x.Window),
		x.On("click", "close($refs.selectroot)", x.Outside),
		x.X("anchor.bottom.offset.6", "$refs.trigger"),
		x.Bind("data-state", "opened ? 'open' : 'closed'"),
		html.Role("listbox"),
		html.Data("slot", "select-content"),
		html.Class("relative z-50 min-w-[8rem] w-full overflow-y-auto overflow-x-hidden rounded-md border bg-popover text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"),
		html.Div(
			components.ItemsOf(
				items,
				html.Class("p-1 h-[var(--radix-select-trigger-height)] w-full min-w-[var(--radix-select-trigger-width)] scroll-my-1"),
			)...,
		),
	)
}

func Item(items ...html.Item) html.Node {
	items = components.ItemsOf(
		items,
		Select,
		x.Bind("aria-selected", "false"),
		x.On("mouseenter", "$focus.focus($el)"),
		x.On("mouseleave", "$el.blur()"),
		x.On("keydown", "$el.click()", "enter"),
		x.On("keydown", "$el.click()", "space"),
		html.Data("slot", "select-item"),
		html.Role("option"),
		html.Class("focus:bg-accent focus:text-accent-foreground [&_svg:not([class*='text-'])]:text-muted-foreground relative flex w-full cursor-default items-center gap-2 rounded-sm py-1.5 pr-8 pl-2 text-sm outline-hidden select-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 *:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2"),
	)

	items = append(items,
		html.Span(
			x.Show("value === $el.parentElement.getAttribute('data-value')"),
			html.Aria("hidden", "true"),
			html.Class("absolute right-2 flex size-3.5 items-center justify-center"),

			lucide.Check(html.Class("h-4 w-4")),
		),
	)
	return html.Div(items...)
}

var Open = OpenOn("click")
var Close = CloseOn("click")
var Select = SelectOn("click")

func OpenOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "open('mouse', $refs.selectroot)")
}

func CloseOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "close($refs.selectroot)")
}

func SelectOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "select($el.getAttribute('data-value'), $refs.selectroot); $nextTick(); close($refs.selectroot)")
}
