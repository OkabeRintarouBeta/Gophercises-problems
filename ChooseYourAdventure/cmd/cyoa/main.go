package main

import (
	"flag"
	"net/http"
	"os"
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
	err = http.ListenAndServe(":4000", app.routes())
	os.Exit(1)
}
