package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	. "github.com/okaberintaroubeta/SitemapBuilder/Link"
)

var visited map[string]bool

func main() {
	siteURL := flag.String("url", "https://www.calhoun.io/", "URL for the target website")
	flag.Parse()
	visited = make(map[string]bool)

	var hrefList []string
	hrefList = append(hrefList, *siteURL)
	var result []string
	result = bfs(hrefList, 2)
	for _, href := range result {
		fmt.Println(href)
	}

}

func getHTML(URL string) (string, error) {
	resp, err := http.Get(URL)
	// fmt.Println("----------------------------")
	if err != nil {
		// fmt.Errorf("GET error: %v", err)
		return "", err
	}
	// fmt.Println("----------------------------")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// fmt.Printf("Status error: %v\n", resp.StatusCode)
		return "", err
	}
	// fmt.Println("----------------------------")

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Printf("Read body: %v\n", err)
		return "", err
	}
	// fmt.Println("----------------------------")
	return string(data), nil
}

func bfs(hrefList []string, depth int) []string {

	var result []string
	for len(hrefList) > 0 && depth > 0 {
		var newHrefList []string
		for _, href := range hrefList {
			result = append(result, href)
			// fmt.Println(len(result))
			content, err := getHTML(href)
			// fmt.Println(content)
			if err != nil {
				fmt.Printf("Content of %s can't be retrieved\n ERROR: %s\n", href, err)
				continue
			}
			r := strings.NewReader(content)
			links, err := ParseLink(r)
			if err != nil && len(links) > 0 {
				continue
			}
			for _, l := range links {
				// fmt.Println(l.Href)
				has := visited[l.Href]
				if has == false && isSameSoure(l.Href, hrefList[0]) {
					visited[l.Href] = true
					hrefList = append(hrefList, l.Href)
					newHrefList = append(newHrefList, l.Href)
				} else if has == false && len(l.Href) > 1 && l.Href[0] == '/' {
					// fmt.Println("added url:", hrefList[0]+l.Href)
					newHref := makeHref(l.Href, hrefList[0])
					if visited[newHref] == false {
						visited[newHref] = true
						hrefList = append(hrefList, newHref)
						newHrefList = append(newHrefList, newHref)
					}
				}
			}
		}
		hrefList = newHrefList
		depth -= 1
	}

	return result
}

func isSameSoure(newURL, source string) bool {
	if len(newURL) < len(source) {
		return false
	}
	return newURL[:len(source)] == source
}

func makeHref(ending, original string) string {
	if ending[0] == '/' && original[len(original)-1] == '/' {
		return original + ending[1:]
	} else if ending[0] != '/' && original[len(original)-1] != '/' {
		return original + "/" + ending
	} else {
		return original + ending
	}
}
