package main

import (
	"fmt"
	"net/http"
	"s_renovation.net/internal/data"
)

func (app *application) createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `form:"name"`
		Email   string `form:"email"`
		Details string `form:"details"`
	}

	err := app.decodePostForm(r, &input)
	fmt.Println(input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	order := &data.Order{
		Name:    input.Name,
		Email:   input.Email,
		Details: input.Details,
	}

	err = app.models.Order.Insert(order)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

//func (app *application) showAllOrders(w http.ResponseWriter, r *http.Request) {
//	orders, err := app.models.Order.GelAll()
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//		return
//	}
//}
