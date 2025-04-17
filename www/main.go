package main

import (
	"log"

	"github.com/canpacis/pacis/pages"
	"github.com/canpacis/pacis/www/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := app.InitDocs(); err != nil {
		log.Fatal(err)
	}
	err := godotenv.Load("www/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pages.Serve(":8080", app.Router(nil), nil)
}
