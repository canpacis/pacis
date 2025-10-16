package server

import (
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/server/metadata"
)

type NotFound struct{}

func (*NotFound) Metadata() *metadata.Metadata {
	return &metadata.Metadata{
		Title: "Page Not Found",
	}
}

func (*NotFound) Page() html.Node {
	return html.Div(
		html.StyleAttr("min-height: 100dvh; display: flex; flex-direction: column; justify-content: center; place-items: center; font-family: system-ui, sans-serif;"),

		html.H1(
			html.StyleAttr("font-size: 2.25rem; font-weight: bold; margin: 1rem 0;"),

			html.Text("404"),
		),
		html.P(
			html.StyleAttr("margin: 0;"),

			html.Text("Page Not Found"),
		),
	)
}

var NotFoundPage = &NotFound{}
