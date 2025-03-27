package components

import "github.com/canpacis/pacis/renderer"

func AvatarFallback(props ...renderer.Renderer) *renderer.Element {
	ps := []renderer.Renderer{
		renderer.Class("bg-muted flex size-full items-center justify-center rounded-full"),
	}
	ps = append(ps, props...)

	return renderer.Div(ps...)
}

func AvatarImage(props ...renderer.Renderer) *renderer.Element {
	ps := []renderer.Renderer{
		renderer.Class("aspect-square size-full"),
	}
	ps = append(ps, props...)

	return renderer.Img(ps...)
}

func Avatar(props ...renderer.Renderer) *renderer.Element {
	ps := []renderer.Renderer{
		renderer.Class("relative flex size-8 shrink-0 overflow-hidden rounded-full"),
	}
	ps = append(ps, props...)

	return renderer.Div(ps...)
}
