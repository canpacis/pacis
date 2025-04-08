package components

import (
	"bytes"
	"context"
	"fmt"

	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

func Checkbox(props ...h.I) h.Element {
	el := h.S(props...)
	_, checked := el.GetAttribute("checked")
	idattr, hasid := el.GetAttribute("id")

	var id string
	if !hasid {
		id = randid()
	} else {
		// TODO: should remove the id from label
		var buf bytes.Buffer
		idattr.Render(context.Background(), &buf)
		id = buf.String()
	}

	// TODO: route the checked property to input inside
	return h.Lbl(
		Join(
			props,
			h.Class("text-sm gap-2 items-center inline-flex w-fit-content cursor-pointer"), h.HtmlFor(id),
			h.Span(
				h.Class("peer border-input dark:bg-input/30 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground dark:data-[state=checked]:bg-primary data-[state=checked]:border-primary focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive size-4 aspect-square shrink-0 rounded-[4px] border shadow-xs transition-shadow outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50"),
				X("data", fmt.Sprintf("checkbox(%t, '%s')", checked, id)),
				h.Attr(":data-state", "checked ? 'checked' : 'unchecked'"),

				h.Span(
					h.Class("items-center justify-center text-current transition-none"),
					h.Attr(":class", "!checked ? 'hidden' : 'flex'"),
					icons.Check(h.Class("size-3.5")),
				),

				h.Inpt(
					h.ID(id),
					h.Type("checkbox"),
					h.Class("sr-only"),
					X("bind:checked", "checked"),
					On("change", "checked = !checked; await $nextTick(); $dispatch('changed', { checked: checked, event: $event });"),
				),
			),
		)...)
}

// Returns an attribte that toggles the related checkbox upon given event
func ToggleCheckboxOn(event string) h.Attribute {
	return h.Attr(event, "toggleCheckbox()")
}

// An attribute that toggles the related checkbox on click
var ToggleCheckbox = ToggleCheckboxOn("click")
