package link

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseLink(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil, err
	}
	var links []Link
	nodes := splitNodes(doc)
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	// for _, link := range links {
	// 	fmt.Println(link.Href, ":", link.Text)
	// }
	return links, nil
}

func buildLink(node *html.Node) Link {
	var link Link
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			break
		}
	}
	link.Text = getText(node)
	return link
}

func getText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	if node.Type != html.ElementNode {
		return ""
	}
	var result string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		result += getText(c)
	}
	return strings.Join(strings.Fields(result), " ")

}

func splitNodes(node *html.Node) [](*html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}
	var nodes []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, splitNodes(c)...)
	}
	return nodes
}
