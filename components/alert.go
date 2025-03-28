package components

import (
	"context"
	"io"

	r "github.com/canpacis/pacis/renderer"
)

func AlertTitle(props ...r.I) r.Element {
	ps := []r.I{
		r.Class("col-start-2 line-clamp-1 min-h-4 font-medium tracking-tight"),
		r.Role("alert-title"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

func AlertDescription(props ...r.I) r.Element {
	ps := []r.I{
		r.Class("text-muted-foreground col-start-2 grid justify-items-start gap-1 text-sm [&_p]:leading-relaxed"),
		r.Role("alert-description"),
	}
	ps = append(ps, props...)

	return r.Div(ps...)
}

type AlertVariant int

const (
	AlertVariantDefault = AlertVariant(iota)
	AlertVariantDestructive
)

func (v AlertVariant) Render(ctx context.Context, w io.Writer) error {
	return r.Class(v.GetValue().(string)).Render(ctx, w)
}

func (AlertVariant) GetKey() string {
	return "class"
}

func (v AlertVariant) GetValue() any {
	switch v {
	case AlertVariantDefault:
		return "bg-card text-card-foreground"
	case AlertVariantDestructive:
		return "text-destructive bg-card [&>svg]:text-current *:data-[slot=alert-description]:text-destructive/90"
	default:
		panic("invalid alert variant property")
	}
}

func Alert(props ...r.I) r.Element {
	var variant AlertVariant

	ps := []r.I{
		r.Class("relative w-full rounded-lg border px-4 py-3 text-sm grid has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr] grid-cols-[0_1fr] has-[>svg]:gap-x-3 gap-y-0.5 items-start [&>svg]:size-4 [&>svg]:translate-y-0.5 [&>svg]:text-current"),
		r.Role("alert"),
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
	return r.Div(ps...)
}
