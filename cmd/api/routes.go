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

	router.HandlerFunc(http.MethodGet, "/", app.home)

	router.HandlerFunc(http.MethodGet, "/signup", app.userSignup)
	router.HandlerFunc(http.MethodGet, "/signin", app.userSignin)

	router.HandlerFunc(http.MethodGet, "/user/profile", app.showUserProfile)
	router.HandlerFunc(http.MethodGet, "/user/settings", app.showUserSettings)
	router.HandlerFunc(http.MethodGet, "/user/orders", app.showUserOrders)

	router.HandlerFunc(http.MethodGet, "/admin", app.showAdminPanel)
	router.HandlerFunc(http.MethodGet, "/admin/customers", app.showAdminPanelCustomers)
	router.HandlerFunc(http.MethodGet, "/admin/orders", app.showAdminPanelOrders)

	return router
}
