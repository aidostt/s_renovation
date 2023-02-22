package main

import (
	"encoding/json"
	"net/http"
	"s_renovation.net/internal/data"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.PrintInfo("fake error", nil) // nigga
	w.Write([]byte("Successful connection"))

}

func (app *application) createFormHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Form
	}

	err := json.NewDecoder(r.Body).Decode(&input)

	form := &data.Form{
		Floor:     input.Floor,
		Plinth:    input.Plinth,
		Door:      input.Door,
		Toilet:    input.Toilet,
		Socket:    input.Socket,
		Plumb:     input.Plumb,
		Paint:     input.Paint,
		Wallpaper: input.Wallpaper,
	}
	err = app.models.Form.Insert(form)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
