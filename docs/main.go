package main

import (
	"net/http"

	"github.com/canpacis/pacis/docs/app"
	"github.com/canpacis/pacis/pages"
	"github.com/canpacis/pacis/ui/components"
	"github.com/canpacis/pacis/ui/html"
)

func main() {
	mux := http.NewServeMux()
	head := components.CreateHead("/ui/")

	mux.Handle("/", pages.NewPageRoute("/", app.HomeLayout, func(pc *pages.PageContext) html.I {
		return html.Div(html.Text("Hello, World!"))
	}))
	mux.Handle("/ui/", http.StripPrefix("/ui/", head.Handler()))
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./docs/public"))))

	http.ListenAndServe(":8080", mux)
}
