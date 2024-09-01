package main

import (
	"net/http"

	"github.com/berberapan/my-stuff/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /signup", app.signup)

	return mux
}
