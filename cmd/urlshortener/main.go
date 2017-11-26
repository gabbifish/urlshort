package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/gabbifish/urlshort/shortener"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := shortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback

	yaml, _ := ioutil.ReadFile("mapping.yaml")
	yamlHandler, err := shortener.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	// http.ListenAndServe(":8080", mapHandler)
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
