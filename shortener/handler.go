package shortener

import (
	"gopkg.in/yaml.v2"
	"net/http"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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
		url, ok := pathsToUrls[slug]
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

	var entries []URLExpand
	err := yaml.Unmarshal(yml, &entries)

	// Turn list of yaml entries into map
	m := make(map[string]string)
	for _, entry := range entries {
		m[entry.Slug] = entry.Full
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Path
		url, ok := m[slug]
		if ok {
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}), err
}

// SQLHandler returns a HandlerFunc that queries a slug's full URL in a SQL
// database (here, SQLlite3 for simplicity). If there is a database error,
// SQLHandler's returned HandlerFunc will respond with a 500 Internal Server
// Error status and an error message. Otherwise, if the slug-URL mapping exists
// in the SQL db, then the HandlerFunc will redirect to the appropriate page.
// If mapping does not exist, SQLHandler passes request to fallback Handler.
func SQLHandler(dbName string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the appropriate mapping from the request
		slug := r.URL.Path
		url, getUrlFromSqlErr := getURLFromDB(dbName, slug)
		if getUrlFromSqlErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Oops! SQL database derped."))
		} else if url == "" {  // If no URL is found, then do the default
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusSeeOther)
		}
	})
}

func getURLFromDB(dbName string, path string) (string, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return "", err
	}

	// Write SQL query
	query := fmt.Sprintf("SELECT url FROM Mappings WHERE slug = '%s'", path)
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}

	// Retrieve result from rows (there should be only one row,
	// so we use if rows.Next(), not for rows.Next()).2
	defer rows.Close()
	url := ""
	if rows.Next() {
		err := rows.Scan(&url)
		if err != nil {
			return url, err
		}
	}
	return url, nil
}
