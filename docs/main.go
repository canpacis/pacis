package main

import (
	"embed"
	"net/http"

	"github.com/canpacis/pacis/docs/app"
	. "github.com/canpacis/pacis/pages"
	"github.com/canpacis/pacis/pages/middleware"
)

//go:embed public
var public embed.FS

func main() {
	router := Routes(
		Public(public, "public"),
		Layout(app.Layout),
		Middleware(middleware.Theme),
		Page(app.HomePage),

		Route(
			Path("docs"),
			Layout(app.DocLayout),

			Route(Path("introduction"), Page(app.Introduction)),
			Route(Path("installation"), Page(app.Installation)),
			Route(Path("components"), Redirect("/docs/alert")),
			Route(Path("alert"), Page(app.AlertPage)),
		),
	)

	http.ListenAndServe(":8080", router.Handler())
}
