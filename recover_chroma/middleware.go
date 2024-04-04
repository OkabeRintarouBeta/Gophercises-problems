package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)

				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "text/html")
				stack := string(debug.Stack())
				lineString := makeLines(strings.Split(stack, "\n"))
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, lineString)
				return
			}
		}()
		app.ServeHTTP(w, r)
	}
}
