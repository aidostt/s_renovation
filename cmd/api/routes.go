package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) Router() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/createForm", app.createFormHandler)
	router.HandlerFunc(http.MethodPost, "/signup", app.registerUserHandler)
	return router
}
