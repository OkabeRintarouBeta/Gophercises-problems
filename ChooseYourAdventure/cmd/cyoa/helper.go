package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	cyoa "github.com/okaberintaroubeta/cyoa/models"
)

func parseJson(filePath string) (cyoa.Story, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	defer jsonFile.Close()
	d := json.NewDecoder(jsonFile)
	var story cyoa.Story
	if err = d.Decode(&story); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	// fmt.Printf("%+v\n", story)
	return story, nil
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Server error: %v", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template does not exist", page)
		app.serverError(w, r, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}
