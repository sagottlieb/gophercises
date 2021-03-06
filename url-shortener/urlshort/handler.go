package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if dest, found := pathsToUrls[r.URL.Path]; found {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
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

	yamlData, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	redirectsMap := buildMap(yamlData)

	return MapHandler(redirectsMap, fallback), nil
}

type redirect struct {
	Path string
	Url  string
}

func parseYAML(yml []byte) ([]redirect, error) {
	results := []redirect{}

	err := yaml.Unmarshal(yml, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func buildMap(redirects []redirect) map[string]string {
	out := map[string]string{}

	for _, r := range redirects {
		out[r.Path] = r.Url
	}

	return out
}
