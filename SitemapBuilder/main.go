package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	. "github.com/okaberintaroubeta/SitemapBuilder/Link"
)

var visited map[string]bool

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	XmlNS string `xml:"xmlns,attr"`
}

func main() {
	siteURL := flag.String("url", "http://google.com", "URL for the target website")
	depth := flag.Int("depth", 1, "Search depth for the URL")
	flag.Parse()

	visited = make(map[string]bool)

	var hrefList []string
	hrefList = append(hrefList, *siteURL)
	var result []string
	result = bfs(hrefList, *depth)

	toXml := urlset{
		XmlNS: xmlns,
	}

	xmlFile, err := os.Create("my-file.xml")
	if err != nil {
		fmt.Println("Error creating XML file: ", err)
		return
	}

	xmlFile.WriteString(xml.Header)
	enc := xml.NewEncoder(xmlFile)
	enc.Indent("", "\t")

	for _, res := range result {
		toXml.Urls = append(toXml.Urls, loc{res})
	}
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
}

func bfs(hrefList []string, depth int) []string {

	// Get the base URL
	resp, err := http.Get(hrefList[0])
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	reqUrl := resp.Request.URL
	defer resp.Body.Close()
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	var result []string

	// BFS Search with maximum depth
	for len(hrefList) > 0 && depth >= 0 {
		var newHrefList []string
		for _, href := range hrefList {

			resp, err := http.Get(href)
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				// fmt.Printf("Status error: %v\n", resp.StatusCode)
				continue
			}

			reqUrl := resp.Request.URL

			// Get all links in <a> tag
			links, err := ParseLink(resp.Body)
			if err != nil && len(links) > 0 {
				continue
			}
			// make sure the link is valid before appending it to the final result
			result = append(result, reqUrl.String())

			for _, l := range links {
				var newHref string
				// If the new href appears to be part of the url
				// add it to the end of the base to make a new url
				if len(l.Href) > 1 && (strings.HasPrefix(l.Href, "/") || !strings.HasPrefix(l.Href, base)) {
					newHref = makeHref(l.Href, base)
				} else {
					newHref = l.Href
				}
				// If the new href belongs to the same site as the base
				// Then append it to the new layer of BFS
				if filterURL(newHref, base) {
					newHrefList = append(newHrefList, newHref)
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

// Check if the new URL has already been visited or share the same prefix as the base
func filterURL(newURL, base string) bool {
	if visited[newURL] {
		return false
	}
	visited[newURL] = true
	if strings.HasPrefix(newURL, base) {
		return true
	}
	return false
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
