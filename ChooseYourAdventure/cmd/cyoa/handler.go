package main

import (
	"log"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	chapter, ok := app.story.Get("intro")
	if !ok {
		app.notFound(w)
		return
	}
	w.Header().Set("Content-Type", "text/html")

	data := templateData{}
	data.Chapter = chapter
	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) chapterView(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	log.Printf("the title is: %s", title)
	chapter, exists := app.story.Get(title)
	if !exists {
		app.notFound(w)
	} else {
		data := templateData{}
		data.Chapter = chapter
		app.render(w, r, http.StatusOK, "view.tmpl", data)
	}
}
