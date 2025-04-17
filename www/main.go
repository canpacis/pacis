package main

import (
	"log"

	"github.com/canpacis/pacis/pages"
	"github.com/canpacis/pacis/www/app"
)

func main() {
	if err := app.InitDocs(); err != nil {
		log.Fatal(err)
	}
	pages.Serve(":8080", app.Router(nil), nil)
}
