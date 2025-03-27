package main

import (
	"net/http"

	. "github.com/canpacis/pacis/components"
	. "github.com/canpacis/pacis/renderer"
)

func main() {
	html := Div(
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
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		html.Render(w)
	})

	http.ListenAndServe(":8080", mux)
}
