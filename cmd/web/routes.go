package main

import (
	"net/http"

	"github.com/berberapan/my-stuff/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	mux.HandleFunc("GET /healthz", getHealthz)

	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave, app.authenticate, csrfProtection)

	mux.Handle("GET /{$}", dynamicMiddleware.ThenFunc(app.getHome))
	mux.Handle("GET /signup", dynamicMiddleware.ThenFunc(app.getSignup))
	mux.Handle("POST /signup", dynamicMiddleware.ThenFunc(app.postSignup))
	mux.Handle("GET /login", dynamicMiddleware.ThenFunc(app.getLogin))
	mux.Handle("POST /login", dynamicMiddleware.ThenFunc(app.postLogin))

	protectedMiddleware := dynamicMiddleware.Append(app.requireAuthentication)

	mux.Handle("POST /logout", protectedMiddleware.ThenFunc(app.postLogout))
	mux.Handle("GET /stuff", protectedMiddleware.ThenFunc(app.getStuff))

	standardMiddleware := alice.New(app.logRequest, app.recoverPanic)

	return standardMiddleware.Then(mux)
}
