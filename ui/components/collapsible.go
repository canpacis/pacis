package components

import (
	"fmt"

	h "github.com/canpacis/pacis/ui/html"
)

/*
	An interactive component which expands/collapses a panel.

Usage:

	Collapsible(
		CollapsibleTrigger(
			Button(Text("Trigger")) // <- Pass a button element for accessiblity
		),
		CollapsibleContent(
			P(Text("Content"))
		),
	)
*/
func Collapsible(props ...h.I) h.Element {
	el := h.Div(props...)

	idattr, hasid := el.GetAttribute("id")
	var id string
	if !hasid {
		id = randid()
	} else {
		id = readattr(idattr)
		el.RemoveAttribute("id")
	}

	el.AddAttribute(X("data", fmt.Sprintf("collapsible(false, '%s')", id)))
	return el
}

// A trigger element for collapsible component
func CollapsibleTrigger(trigger h.Element) h.Element {
	trigger.AddAttribute(ToggleCollapsible)
	return trigger
}

// Content slot for the collapsible component
func CollapsibleContent(content h.Element) h.Element {
	content.AddAttribute(X("show", "open"))
	return content
}

/*
	Returns an attribute that toggles the related collapsible upon given event

Usage:

	Collapsible(
		CollapsibleTrigger( ... ),
		CollapsibleContent(
			Button( // <- Another button for toggling
				ToggleCollapsibleOn("mouseover") // <- will toggle the collapsible when hovered

				Text("Toggle outside")
			),
		),
	)
*/
func ToggleCollapsibleOn(event string) h.Attribute {
	return On(event, "toggleCollapsible()")
}

/*
	An attribute that toggles the related collapsible on click

Usage:

	Collapsible(
		CollapsibleTrigger( ... ),
		CollapsibleContent(
			Button( // <- Another button for toggling
				ToggleCollapsible

				Text("Toggle outside")
			),
		),
	)
*/
var ToggleCollapsible = ToggleCollapsibleOn("click")
