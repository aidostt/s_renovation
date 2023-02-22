package main

import "net/http"

func (app *application) showAdminPanel(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "adminPanel.htm", data, r)
}

func (app *application) showAdminPanelCustomers(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "adminPanel_customers.htm", data, r)
}

func (app *application) showAdminPanelOrders(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "adminPanel_orders.htm", data, r)
}
