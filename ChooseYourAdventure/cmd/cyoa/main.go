package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"

	cyoa "github.com/okaberintaroubeta/cyoa/models"
)

type application struct {
	story         *cyoa.Story
	templateCache map[string]*template.Template
}

func main() {
	filePath := flag.String("path", "gopher.json", "Path of the original story")
	flag.Parse()

	templateCache, err := newTemplateCache()
	story, err := parseJson(*filePath)
	if err != nil {
		return
	}
	app := &application{
		story:         &story,
		templateCache: templateCache,
	}
	log.Fatal(http.ListenAndServe(":4000", app.routes()))
}
