package components

import (
	"context"
	"io"

	h "github.com/canpacis/pacis/ui/html"
)

func AlertTitle(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("col-start-2 line-clamp-1 min-h-4 font-medium tracking-tight"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

func AlertDescription(props ...h.I) h.Element {
	ps := []h.I{
		h.Class("text-muted-foreground col-start-2 grid justify-items-start gap-1 text-sm [&_p]:leading-relaxed"),
	}
	ps = append(ps, props...)

	return h.Div(ps...)
}

type AlertVariant int

const (
	AlertVariantDefault = AlertVariant(iota)
	AlertVariantDestructive
)

func (v AlertVariant) Render(ctx context.Context, w io.Writer) error {
	var value string
	switch v {
	case AlertVariantDefault:
		value = "bg-card text-card-foreground"
	case AlertVariantDestructive:
		value = "text-destructive bg-card [&>svg]:text-current *:data-[slot=alert-description]:text-destructive/90"
	default:
		panic("invalid alert variant property")
	}

	return h.Class(value).Render(ctx, w)
}

func (AlertVariant) GetKey() string {
	return "class"
}

func (v AlertVariant) IsEmpty() bool {
	return false
}

func Alert(props ...h.I) h.Element {
	var variant AlertVariant

	ps := []h.I{
		h.Class("relative w-full rounded-lg border px-4 py-3 text-sm grid has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr] grid-cols-[0_1fr] has-[>svg]:gap-x-3 gap-y-0.5 items-start [&>svg]:size-4 [&>svg]:translate-y-0.5 [&>svg]:text-current"),
		h.Role("alert"),
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
	return h.Div(ps...)
}
