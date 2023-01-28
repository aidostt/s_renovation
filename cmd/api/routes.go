package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) Router() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthCheckHandler)

	return router
}
