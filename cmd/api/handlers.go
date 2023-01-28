package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.errorLog.Print("fake error")
	w.Write([]byte("Successful connection"))
}
