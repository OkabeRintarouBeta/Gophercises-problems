package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Query().Get("path")
	lineNum, convErr := strconv.Atoi(r.URL.Query().Get("line"))
	// path, _ := os.Getwd()
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	lexer := lexers.Get("go")
	iterator, err := lexer.Tokenise(nil, b.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	var formatter *html.Formatter
	if convErr == nil { // If there's no error in converting lineNum
		formatter = html.New(html.WithLineNumbers(true), html.HighlightLines([][2]int{{lineNum, lineNum}}))
	} else {
		formatter = html.New(html.WithLineNumbers(true))
	}
	w.Header().Set("Content-Type", "text/html")
	err = formatter.Format(w, style, iterator)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
