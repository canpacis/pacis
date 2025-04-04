package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/canpacis/pacis/docs/app"
	. "github.com/canpacis/pacis/pages"
)

//go:embed public
var public embed.FS

func Theme(ctx Context) Context {
	fmt.Println(ctx)
	return ctx
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println(r.URL)
	// 	cookie, err := r.Cookie("pacis_color_scheme")
	// 	if err != nil {
	// 		switch cookie.Value {
	// 		case "light":
	// 		case "dark":
	// 		default:
	// 			// delete cookie
	// 		}
	// 	}
	// 	h.ServeHTTP(w, r)
	// })
}

func main() {
	router := Routes(
		Public(public, "public"),
		Layout(app.Layout),
		Page(app.HomePage),

		Route(
			Path("docs"),
			Layout(app.DocLayout),

			Route(Path("introduction"), Page(app.Introduction), Middleware(Theme)),
			Route(Path("installation"), Page(app.Installation)),
			Route(Path("components"), Redirect("/docs/alert")),
			Route(Path("alert"), Page(app.AlertPage)),
		),
	)

	http.ListenAndServe(":8080", router.Handler())
}
