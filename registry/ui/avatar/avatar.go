package avatar

import (
	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/x"
)

func New(items ...html.Item) html.Node {
	type Data struct {
		Error *string `json:"error"`
	}

	return html.Span(
		components.ItemsOf(
			items,
			x.Data(Data{}),
			html.Data("slot", "avatar"),
			html.Class("relative flex h-10 w-10 shrink-0 overflow-hidden rounded-full"),
		)...,
	)
}

func Image(items ...html.Item) html.Node {
	return html.Img(
		components.ItemsOf(
			items,
			html.Data("slot", "avatar-image"),
			html.Class("aspect-square h-full w-full"),
			html.Attr("x-show", "error === null"),
			html.Attr("x-on:error", "error = 'failed to load image'"),
		)...,
	)
}
func Fallback(items ...html.Item) html.Node {
	return html.Span(
		components.ItemsOf(
			items,
			html.Data("slot", "avatar-fallback"),
			html.Class("flex h-full w-full items-center justify-center rounded-full bg-muted"),
			html.Attr("x-show", "error !== null"),
		)...,
	)
}
