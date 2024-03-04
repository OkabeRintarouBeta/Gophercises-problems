package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/okaberintaroubeta/Gophercises-problems/urlshort"
)

func main() {

	filePath := flag.String("path", "", "Path to parsed file")
	flag.Parse()
	content, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Println("Error:", err)
	}
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/final
	// `
	var fileHandler http.HandlerFunc

	if (*filePath)[len(*filePath)-4:] == "json" {
		fileHandler, err = urlshort.JSONHandler([]byte(content), mapHandler)
	} else {
		fileHandler, err = urlshort.YAMLHandler([]byte(content), mapHandler)
	}

	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", fileHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
