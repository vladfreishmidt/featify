package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/vladfreishmidt/featify/internal/models"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type templateData struct {
	CurrentYear     int
	Project         *models.Project
	Projects        []*models.Project
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func newTemplateCahce() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
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

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}
