package main

import (
	"log"

	"github.com/canpacis/pacis/pages"
	"github.com/canpacis/pacis/www/app"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("www/.env")
	if err := app.InitDocs(); err != nil {
		log.Fatal(err)
	}
	app.InitAuth()

	pages.Serve(":8080", app.Router(nil), nil)
}
