package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/canpacis/pacis/docs/app"
	p "github.com/canpacis/pacis/pages"
	"github.com/canpacis/pacis/pages/i18n"
	"github.com/canpacis/pacis/pages/middleware"
	"golang.org/x/text/language"
)

//go:embed public
var public embed.FS

//go:embed messages
var messages embed.FS

//go:embed app/markup
var markup embed.FS

//go:embed app/robots.txt
var robots []byte

//go:embed app/sitemap.xml
var sitemap []byte

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type docitem struct {
	path   string
	markup string
}

func main() {
	bundle, err := i18n.Setup(messages, language.English)
	if err != nil {
		log.Fatal(err)
	}

	docs := []docitem{
		{"introduction", "app/markup/introduction.md"},
		{"installation", "app/markup/installation.md"},
		{"alert", "app/markup/alert.md"},
		{"avatar", "app/markup/avatar.md"},
		{"badge", "app/markup/badge.md"},
		{"button", "app/markup/button.md"},
		{"card", "app/markup/card.md"},
		{"checkbox", "app/markup/checkbox.md"},
		{"collapsible", "app/markup/collapsible.md"},
		{"dialog", "app/markup/dialog.md"},
		{"dropdown", "app/markup/dropdown.md"},
		{"input", "app/markup/input.md"},
		{"label", "app/markup/label.md"},
		{"select", "app/markup/select.md"},
		{"tabs", "app/markup/tabs.md"},
	}

	router := p.Routes(
		p.Public(public, "public"),
		p.Layout(app.Layout),
		p.Middleware(middleware.Theme),
		p.Middleware(middleware.Locale(bundle, language.English)),
		p.Page(app.HomePage),

		p.Route(
			p.Path("docs"),
			p.Layout(app.DocLayout),

			p.Route(p.Path("components"), p.Redirect("/docs/alert")),
			p.Map(docs, func(doc docitem, i int) p.RouteItem {
				return p.Route(p.Path(doc.path), p.Page(app.MarkupPage(markup, doc.markup)))
			}),
		),

		p.Robots(robots),
		p.Sitemap(sitemap),
	)

	http.ListenAndServe(":"+getEnv("PORT", "8080"), router.Handler())
}
