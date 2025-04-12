package main

import (
	"net/http"

	"github.com/canpacis/pacis/www/app"
)

func main() {
	router := app.Router()
	http.ListenAndServe(":8081", router)
}
