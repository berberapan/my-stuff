package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.tmpl")
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "signup.tmpl")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
