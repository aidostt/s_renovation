package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"s_renovation.net/internal/data"
	"s_renovation.net/validator"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `form:"name"`
		Surname  string `form:"surname"`
		Phone    string `form:"phone"`
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	err := app.decodePostForm(r, &input)
	fmt.Println(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		CreatedAt: time.Now(),
		Surname:   input.Surname,
		Phone:     input.Phone,
		Email:     input.Email,
		Role:      0,
	}

	err = user.Password.Set(input.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.logger.PrintInfo("failed to validate user", nil)
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.User.Insert(user)
	if err != nil {
		fmt.Printf("\nRegister user handler: %v", err)
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			fmt.Println("i was in default")
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	http.Redirect(w, r, "/signin", http.StatusSeeOther)
	//err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
}

func (app *application) registerUserHandlerGET(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "signup.htm", nil, r)
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	var email string
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateEmail(v, email); !v.Valid() {
		app.logger.PrintInfo("couldn't validate email", nil)
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.User.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	w.Write([]byte(user.Name))
}

func (app *application) userSigninPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `form:"email"`

		Password string `form:"password"`
	}
	err := app.decodePostForm(r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	//validation

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
	}
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	w.Write([]byte("im ok"))
}

func (app *application) userSignin(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "signin.htm", nil, r)
}

func (app *application) showUserProfile(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "userProfile.htm", nil, r)
}

func (app *application) showUserSettings(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "userSettings.htm", nil, r)
}

func (app *application) showUserOrders(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "userOrders.htm", nil, r)
}

func (app *application) showAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.User.GelAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = users
	app.render(w, http.StatusOK, "adminPanel_customers.htm", data, r)
}
