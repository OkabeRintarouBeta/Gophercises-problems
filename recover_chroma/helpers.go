package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func findPattern(stk, pattern string) []int {
	r := regexp.MustCompile(pattern)
	urlIdx := r.FindStringIndex(stk)
	return urlIdx
}

func makeLines(lines []string) string {
	for i, line := range lines {
		urlPattern := `/\S+:`
		urlIndex := findPattern(line, urlPattern)
		if urlIndex != nil {
			lineIdx := findPattern(line, `:\d+`)
			// fmt.Printf("%s\n", line[urlIndex[0]:urlIndex[1]])
			v := url.Values{}
			v.Add("path", line[urlIndex[0]:urlIndex[1]-1])
			v.Add("line", line[lineIdx[0]+1:lineIdx[1]])
			lines[i] = fmt.Sprintf("%s<a href=/debug?%s>%s</a> :%s",
				line[:urlIndex[0]], v.Encode(), line[urlIndex[0]:urlIndex[1]-1], line[urlIndex[1]:])
		}
	}
	return strings.Join(lines, "\n")
}
