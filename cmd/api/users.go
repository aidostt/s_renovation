package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"s_renovation.net/internal/data"
	"s_renovation.net/validator"
	"time"
)

type userSignupForm struct {
	Name                string `form:"name"`
	Surname             string `form:"surname"`
	Phone               string `form:"phone"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `form:"name"`
		Surname  string `form:"surname"`
		Email    string `form:"email"`
		Password string `form:"password"`
	}
	err := app.decodePostForm(r, &input)
	//err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	user := &data.User{
		Name:      input.Name,
		CreatedAt: time.Now(),
		Surname:   input.Surname,
		Email:     input.Email,
	}
	err = user.Password.Set(input.Password, 12)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.errorLog.Print("failed to validate user")
		app.errorLog.Println(v.Errors)
	}
	err = app.models.User.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	var email string
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		app.errorLog.Print("server error response")
	}
	v := validator.New()
	if data.ValidateEmail(v, email); !v.Valid() {
		app.errorLog.Print("couldn't validate email")
		app.errorLog.Println(v.Errors)
	}
	user, err := app.models.User.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	w.Write([]byte(user.Name))
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.htm", data)
}

func (app *application) userSignin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "signin.htm", data)
}

func (app *application) showUserProfile(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "userProfile.htm", data)
}

func (app *application) showUserSettings(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "userSettings.htm", data)
}

func (app *application) showUserOrders(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "userOrders.htm", data)
}
