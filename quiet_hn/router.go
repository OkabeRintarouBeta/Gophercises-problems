package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes(numStories int, tmpl *template.Template) http.Handler {
	router := httprouter.New()
	router.Handler(http.MethodGet, "/", app.mainHandler(numStories, tmpl))
	return router
}
