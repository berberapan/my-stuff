package main

import (
	"fmt"
	"net/http"
	"os"
)

type signupForm struct {
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm"`
}

func (app *application) getHome(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.tmpl")
}

func (app *application) getSignup(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "signup.tmpl")
}

func (app *application) postSignup(w http.ResponseWriter, r *http.Request) {
	var form signupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		fmt.Println("fel 1")
		// TODO fix proper error handling
		os.Exit(1)
	}

	err = app.users.Insert(form.Email, form.Password)
	if err != nil {
		fmt.Println("fel 2")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getHealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
