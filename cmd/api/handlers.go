package main

import (
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	index := template.Must(template.ParseFiles("D:/Nurkuisa/My/University/5 TRIMESTER/Advanced Databases (NoSQL)/s-renovation/ui/signup.htm"))

	index.Execute(w, nil)

}
