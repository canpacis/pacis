package main

import (
	"net/http"
	"os"

	"github.com/canpacis/pacis/www/app"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	router := app.Router()
	http.ListenAndServe(":"+getEnv("PORT", "8080"), router)
}
