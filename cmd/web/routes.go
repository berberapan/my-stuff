package main

import (
	"net/http"

	"github.com/berberapan/my-stuff/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	mux.HandleFunc("GET /healthz", healthz)

	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamicMiddleware.ThenFunc(app.home))
	mux.Handle("GET /signup", dynamicMiddleware.ThenFunc(app.signup))

	standardMiddleware := alice.New(app.logRequest, app.recoverPanic)

	return standardMiddleware.Then(mux)
}
