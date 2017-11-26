package shortener

import (
	"net/http"
	"gopkg.in/yaml.v2"
	"fmt"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Path
		url, ok := pathsToUrls[slug];
		if ok {
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
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
	type URLExpand struct {
    Slug string `yaml:"path"`
    Full string `yaml:"url"`
	}
	type Config struct {
    Entry []URLExpand `yaml:"config"`
	}

	var c Config
	err := yaml.Unmarshal(yml, &c)

	fmt.Printf("--- m:\n%v\n\n", c)

  // Turn list of yaml entries into maps
	m := make(map[string]string)
	for _, entry := range c.Entry {
		m[entry.Slug] = entry.Full
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Path
		url, ok := m[slug];
		if ok {
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}), err
}
