package main

import (
	"embed"
	"net/http"

	"github.com/canpacis/pacis/docs/app"
	p "github.com/canpacis/pacis/pages"
	"github.com/canpacis/pacis/pages/middleware"
)

//go:embed public
var public embed.FS

type docitem struct {
	path   string
	markup string
}

func main() {
	docs := []docitem{
		{"introduction", "./docs/app/markup/introduction.md"},
		{"installation", "./docs/app/markup/installation.md"},
		{"alert", "./docs/app/markup/alert.md"},
		{"avatar", "./docs/app/markup/avatar.md"},
		{"badge", "./docs/app/markup/badge.md"},
		{"button", "./docs/app/markup/button.md"},
		{"card", "./docs/app/markup/card.md"},
		{"checkbox", "./docs/app/markup/checkbox.md"},
		{"collapsible", "./docs/app/markup/collapsible.md"},
		{"dialog", "./docs/app/markup/dialog.md"},
		{"dropdown", "./docs/app/markup/dropdown.md"},
		{"input", "./docs/app/markup/input.md"},
		{"label", "./docs/app/markup/label.md"},
	}

	router := p.Routes(
		p.Public(public, "public"),
		p.Layout(app.Layout),
		p.Middleware(middleware.Theme),
		p.Page(app.HomePage),

		p.Route(
			p.Path("docs"),
			p.Layout(app.DocLayout),

			p.Route(p.Path("components"), p.Redirect("/docs/alert")),
			p.Map(docs, func(doc docitem, i int) p.RouteItem {
				return p.Route(p.Path(doc.path), p.Page(app.MarkupPage(doc.markup)))
			}),
		),
	)

	http.ListenAndServe(":8080", router.Handler())
}
