package main

import (
	"path/filepath"
	"text/template"

	cyoa "github.com/okaberintaroubeta/cyoa/models"
)

func add(a, b int) int {
	return a + b
}

type templateData struct {
	Chapter cyoa.Chapter
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	funcMap := template.FuncMap{
		"add": add,
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Parse the base template file into a template set
		ts, err := template.New(name).Funcs(funcMap).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
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
