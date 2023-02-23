package main

import (
	"fmt"
	"net/http"
	"s_renovation.net/internal/data"
)

func (app *application) createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name       string `form:"name"`
		Phone      string `form:"phone"`
		Pack       string `form:"pack"`
		Additional bool   `form:"additional"`
		Details    string `form:"details"`
	}

	err := app.decodePostForm(r, &input)
	fmt.Println(input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	order := &data.Order{
		Name:       input.Name,
		Phone:      input.Phone,
		Pack:       input.Pack,
		Additional: input.Additional,
		Details:    input.Details,
	}

	err = app.models.Order.Insert(order)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//func (app *application) showAllOrders(w http.ResponseWriter, r *http.Request) {
//	orders, err := app.models.Order.GelAll()
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//		return
//	}
//}
