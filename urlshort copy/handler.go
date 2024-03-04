package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
type jsonPathUrl struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}
type PathUrl interface{}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}

	// Convert to slice of PathUrl interface
	var interfaceSlice []PathUrl = make([]PathUrl, len(pathUrls))
	for i, d := range pathUrls {
		interfaceSlice[i] = d
	}

	pathToUrls := buildMap(interfaceSlice)
	return MapHandler(pathToUrls, fallback), nil
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []jsonPathUrl
	err := json.Unmarshal(jsn, &pathUrls)
	if err != nil {
		return nil, err
	}

	// Convert to slice of PathUrl interface
	var interfaceSlice []PathUrl = make([]PathUrl, len(pathUrls))
	for i, d := range pathUrls {
		interfaceSlice[i] = d
	}

	pathToUrls := buildMap(interfaceSlice)
	return MapHandler(pathToUrls, fallback), nil
}

func buildMap(pathUrls []PathUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		switch pu := pu.(type) {
		case pathUrl:
			pathToUrls[pu.Path] = pu.URL
		case jsonPathUrl:
			pathToUrls[pu.Path] = pu.URL
		default:
			// handle error or unexpected type
		}
	}
	return pathToUrls
}
