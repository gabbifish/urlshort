package main

import (
	"fmt"
	"net/http"

	short "github.com/gabbifish/urlshort/shortener"
)

func main() {
	mux := defaultMux()

	// Build the mapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapSlugToUrl := short.NewMapSlugToURL(pathsToUrls)
	mapHandler := short.NewHandlerFromSlugToURLClient(mapSlugToUrl, mux)

	// Build the yamlHandler using the mapHandler as the fallback
	yamlSlugToUrl, yamlSlugToUrlErr := short.NewYAMLSlugToURL("mapping.yaml")
	if yamlSlugToUrlErr != nil {
		fmt.Println("yamlHandler init failed:", yamlSlugToUrlErr)
	}
	yamlHandler := short.NewHandlerFromSlugToURLClient(yamlSlugToUrl, mapHandler)

	// Build the sqlHandler using the yamlHandler as the fallback
	sqlSlugToUrl, sqlSlugToUrlErr := short.NewSQLSlugToURL("mapping.db")
	if sqlSlugToUrlErr != nil {
		fmt.Println("sqlHandler init failed:", sqlSlugToUrlErr)
	}
	sqlHandler := short.NewHandlerFromSlugToURLClient(sqlSlugToUrl, yamlHandler)

	// Start server
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", sqlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
