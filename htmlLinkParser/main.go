package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	link "github.com/okaberintaroubeta/htmlLinkParser/Link"
)

func main() {
	filepath := flag.String("path", "examples/ex1.html", "Path to the HTML file to be parsed")
	flag.Parse()
	fileContent, err := os.ReadFile(*filepath)
	fmt.Println(*filepath)
	if err != nil {
		fmt.Printf("Error:", err)
	}
	r := strings.NewReader(string(fileContent))
	link.ParseLink(r)

}
