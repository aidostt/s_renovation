package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"s_renovation.net/internal/data"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     *data.Form
	Snippets    []*data.Form
	Form        any
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.htm")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles("./ui/html/base.htm")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
