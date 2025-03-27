package main

import (
	"fmt"
	"net/http"

	. "github.com/canpacis/pacis/components"
	. "github.com/canpacis/pacis/renderer"
)

func main() {
	html := Html(
		Head(
			Script(Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4")),
			Script(Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js")),
		),
		Body(
			Div(
				Avatar(
					AvatarImage(Src("https://cloud.appwrite.io/v1/storage/buckets/67dc16070005054cb3c3/files/67dc16220011a2637b55/view?project=ksalt&mode=admin")),
					AvatarFallback(Text("MC")),
				),
				Alert(
					AlertTitle(
						Text("This is alert title"),
					),
					AlertDescription(
						Text("This is the description."),
					),
				),
			),
			Div(
				Attr("x-data", "{ count: 0 }"),

				Button(
					On("click", "count++"),

					Text("Increment"),
				),
				Span(Attr("x-text", "count")),
			),
		),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		fmt.Println(html.Render(w))
	})

	http.ListenAndServe(":8080", mux)
}
