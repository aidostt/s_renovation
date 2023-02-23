package main

import (
	"net/http"
	"s_renovation.net/internal/data"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	d := app.newTemplateData(r)
	d.Form = &data.User{}
	app.render(w, http.StatusOK, "index.htm", d, r)
}
