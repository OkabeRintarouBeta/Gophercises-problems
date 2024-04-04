package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/panic", panicDemo)
	r.HandleFunc("/panic-after", panicAfterDemo)

	r.HandleFunc("/debug", sourceCodeHandler)
	r.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", devMw(r)))
}
