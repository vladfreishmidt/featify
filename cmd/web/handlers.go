package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vladfreishmidt/featify/internal/models"
)

// dashboard handler.
func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	projects, err := app.projects.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Projects = projects

	app.render(w, http.StatusOK, "dashboard.tmpl.html", data)
}

// projectView handler.
func (app *application) projectView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	project, err := app.projects.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, http.StatusOK, "project-view.tmpl.html", &templateData{Project: project})
}

// projectCreate handler.
func (app *application) projectCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	name := "Google Cloud"
	description := "This is a test description for the Google Cloud project"

	id, err := app.projects.Insert(name, description)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/project/view?id=%d", id), http.StatusSeeOther)
}

// projectList handler.
func (app *application) projectList(w http.ResponseWriter, r *http.Request) {
	projects, err := app.projects.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "project-list.tmpl.html", &templateData{Projects: projects})
}
