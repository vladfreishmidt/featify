package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.dashboard)
	router.HandlerFunc(http.MethodGet, "/project/view/:id", app.projectView)
	router.HandlerFunc(http.MethodGet, "/project/create", app.projectCreate)
	router.HandlerFunc(http.MethodPost, "/project/create", app.projectCreatePost)
	router.HandlerFunc(http.MethodGet, "/projects", app.projectList)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
