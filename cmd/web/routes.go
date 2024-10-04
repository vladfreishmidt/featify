package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.dashboard)
	mux.HandleFunc("/project/view", app.projectView)
	mux.HandleFunc("/project/create", app.projectCreate)
	mux.HandleFunc("/projects", app.projectList)

	return mux
}
