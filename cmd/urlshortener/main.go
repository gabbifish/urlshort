package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	short "github.com/gabbifish/urlshort/shortener"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := short.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback

	yaml, _ := ioutil.ReadFile("mapping.yaml")
	yamlHandler, err := short.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	sqlHandler := short.SQLHandler("./mapping.db", yamlHandler)
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
