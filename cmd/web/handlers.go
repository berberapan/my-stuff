package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/berberapan/my-stuff/internal/models"
	"github.com/berberapan/my-stuff/internal/validator"
)

type signupForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	ConfirmPassword     string `form:"confirm"`
	validator.Validator `form:"-"`
}

type loginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type addItemForm struct {
	Name                string `form:"name"`
	Description         string `form:"description,omitempty"`
	Accessories         string `form:"accessories,omitempty"`
	Place               string `form:"place,omitempty"`
	AdditionalNotes     string `form:"additional_notes,omitempty"`
	validator.Validator `form:"-"`
}

func (app *application) getHome(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) getSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = signupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) postSignup(w http.ResponseWriter, r *http.Request) {
	var form signupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank.")
	form.CheckField(validator.MatchesRequiredPattern(form.Email, validator.ValidEmailRegex), "email", "This field must be a valid email address.")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank.")
	form.CheckField(validator.MinChars(form.Password, 10), "password", "This field must be at least 10 characters long.")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirm", "This field cannot be blank.")
	form.CheckField(validator.MatchesOtherPassword(form.Password, form.ConfirmPassword), "confirm", "This field needs to match the password field.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Email, form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) getStuff(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	items, err := app.items.AllItems(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Items = items
	app.render(w, r, http.StatusOK, "stuff.tmpl", data)
}

func (app *application) getItem(w http.ResponseWriter, r *http.Request) {
	urlID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	items, err := app.items.AllItems(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	var item models.Item
	found := false
	for _, i := range items {
		if i.ID == urlID {
			item = i
			found = true
			break
		}
	}
	if !found {
		http.NotFound(w, r)
		return
	}
	data := app.newTemplateData(r)
	data.Item = item
	app.render(w, r, http.StatusOK, "item.tmpl", data)
}

func (app *application) getLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) postLogin(w http.ResponseWriter, r *http.Request) {
	var form loginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank.")
	form.CheckField(validator.MatchesRequiredPattern(form.Email, validator.ValidEmailRegex), "email", "This field must a valid email address.")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect.")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/stuff", http.StatusSeeOther)
}

func (app *application) postLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) getAddStuff(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = addItemForm{}
	app.render(w, r, http.StatusOK, "addstuff.tmpl", data)
}

func (app *application) postAddStuff(w http.ResponseWriter, r *http.Request) {
	var form addItemForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank.")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "addstuff.tmpl", data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err = app.items.Insert(form.Name, form.Description, form.Accessories, form.Place, form.AdditionalNotes, id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/stuff", http.StatusSeeOther)
}

func getHealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
