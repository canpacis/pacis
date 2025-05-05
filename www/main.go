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
	if err := app.Init(); err != nil {
		log.Fatal(err)
	}
	router, err := app.Router(nil)
	if err != nil {
		log.Fatal(err)
	}

	pages.Serve(":8081", router, nil)
}
