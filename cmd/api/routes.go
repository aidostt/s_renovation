package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) Router() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/createForm", app.createFormHandler)
	router.HandlerFunc(http.MethodPost, "/signup", app.registerUserHandler)

	router.HandlerFunc(http.MethodGet, "/index", app.home)

	return router
}
