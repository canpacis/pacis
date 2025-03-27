package components

import (
	"io"

	"github.com/canpacis/pacis/renderer"
)

func AlertTitle(props ...renderer.Renderer) *renderer.Element {
	ps := []renderer.Renderer{
		renderer.Class("col-start-2 line-clamp-1 min-h-4 font-medium tracking-tight"),
		renderer.Role("alert-title"),
	}
	ps = append(ps, props...)

	return renderer.Div(ps...)
}

func AlertDescription(props ...renderer.Renderer) *renderer.Element {
	ps := []renderer.Renderer{
		renderer.Class("text-muted-foreground col-start-2 grid justify-items-start gap-1 text-sm [&_p]:leading-relaxed"),
		renderer.Role("alert-description"),
	}
	ps = append(ps, props...)

	return renderer.Div(ps...)
}

type AlertVariant int

const (
	AlertVariantDefault = AlertVariant(iota)
	AlertVariantDestructive
)

func (v AlertVariant) Render(w io.Writer) error {
	class := ""

	switch v {
	case AlertVariantDefault:
		class = "bg-card text-card-foreground"
	case AlertVariantDestructive:
		class = "text-destructive bg-card [&>svg]:text-current *:data-[slot=alert-description]:text-destructive/90"
	default:
		panic("invalid alert variant property")
	}

	return renderer.Class(class).Render(w)
}

func (AlertVariant) Key() string {
	return "class"
}

func Alert(props ...renderer.Renderer) *renderer.Element {
	var variant AlertVariant

	ps := []renderer.Renderer{
		renderer.Class("relative w-full rounded-lg border px-4 py-3 text-sm grid has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr] grid-cols-[0_1fr] has-[>svg]:gap-x-3 gap-y-0.5 items-start [&>svg]:size-4 [&>svg]:translate-y-0.5 [&>svg]:text-current"),
		renderer.Role("alert"),
	}

	for _, prop := range props {
		switch prop := prop.(type) {
		case AlertVariant:
			variant = prop
		default:
			ps = append(ps, prop)
		}
	}

	ps = append(ps, variant)
	return renderer.Div(ps...)
}
