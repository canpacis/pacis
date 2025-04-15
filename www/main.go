package main

import (
	"log"
	"os"

	"github.com/canpacis/pacis/pages"
	"github.com/canpacis/pacis/www/app"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	if err := app.InitDocs(); err != nil {
		log.Fatal(err)
	}
	pages.Serve(":"+getEnv("PORT", "8080"), app.Router(nil))
}
