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

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	// unprotected
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	protected := dynamic.Append(app.requireAuthentication)

	// authenticated-only
	router.Handler(http.MethodGet, "/", protected.ThenFunc(app.dashboard))
	router.Handler(http.MethodGet, "/project/view/:id", protected.ThenFunc(app.projectView))
	router.Handler(http.MethodGet, "/project/create", protected.ThenFunc(app.projectCreate))
	router.Handler(http.MethodPost, "/project/create", protected.ThenFunc(app.projectCreatePost))
	router.Handler(http.MethodGet, "/projects", protected.ThenFunc(app.projectList))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
