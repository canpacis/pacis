package components

import (
	"time"

	h "github.com/canpacis/pacis/ui/html"
)

func Tooltip(content string, delay time.Duration, trigger h.Element, props ...h.I) h.Node {
	trigger.AddAttribute(X("ref", "anchor"))

	return h.Div(
		X("data", "tooltip"),
		QueueOpenTooltip(delay),
		AbortTooltip,

		h.Span(
			Join(
				props,
				X("cloak"),
				X("show", "open"),
				X("transition:enter-start", "scale-90 opacity-0"),
				X("transition:enter-end", "scale-100 opacity-100"),
				X("transition:leave-start", "scale-100 opacity-100"),
				X("transition:leave-end", "scale-90 opacity-0"),
				Anchor(VTop, HCenter, 8),
				h.Class("z-50 overflow-hidden rounded-md bg-primary px-3 py-1.5 text-xs text-primary-foreground transition ease-in-out duration-100 pointer-events-none"),

				h.Text(content),
			)...,
		),
		trigger,
	)
}

func OpenTooltipOn(event string) h.Attribute {
	return On(event, "openTooltip()")
}

var OpenTooltip = OpenTooltipOn("mouseenter")

func CloseTooltipOn(event string) h.Attribute {
	return On(event, "closeTooltip()")
}

var CloseTooltip = CloseTooltipOn("mouseleave")

func QueueOpenTooltipOn(event string, delay time.Duration) h.Attribute {
	return On(event, fn("queueOpenTooltip", int(delay.Milliseconds())))
}

func QueueOpenTooltip(delay time.Duration) h.Attribute {
	return QueueOpenTooltipOn("mouseenter", delay)
}

func AbortTooltipOn(event string) h.Attribute {
	return On(event, "abortTooltip()")
}

var AbortTooltip = AbortTooltipOn("mouseleave")
