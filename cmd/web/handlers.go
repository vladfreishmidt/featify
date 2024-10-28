package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/vladfreishmidt/featify/internal/models"
	"github.com/vladfreishmidt/featify/internal/validator"
)

type projectCreateForm struct {
	Name        string
	Description string
	validator.Validator
}

// dashboard handler.
func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
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
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

// projectCreatePost handler.
func (app *application) projectCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := projectCreateForm{
		Name:        r.PostForm.Get("name"),
		Description: r.PostForm.Get("description"),
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 100), "name", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "project-create.tmpl.html", data)
		return
	}

	id, err := app.projects.Insert(form.Name, form.Description)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/project/view/%d", id), http.StatusSeeOther)
}

// projectCreate handler.
func (app *application) projectCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = projectCreateForm{}

	app.render(w, http.StatusOK, "project-create.tmpl.html", data)
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
