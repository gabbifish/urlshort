package shortener

import (
	"testing"
	"net/http"
	"fmt"
	"net/http/httptest"
	"io/ioutil"
	"strings"
)

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

const HelloResponseBody = "Hello, world!"

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, HelloResponseBody)
}

func TestMalformedYAML(t *testing.T) {
	mux := defaultMux()

	yaml := `
    - path: /urlshort
    url: https://github.com/gophercises/urlshort
    -
    url: https://github.com/gophercises/urlshort/tree/solution`

    _, err := YAMLHandler([]byte(yaml), mux)
    if err == nil {
    	t.Fatal("Incorrectly formed YAML file accepted.")
	}
}

func TestMapHandlerFallback(t *testing.T) {
	mux := defaultMux()
	mapHandler := MapHandler(map[string]string{}, mux)

	// Come up with an HTTP request
	httpRequest, newHttpRequestErr := http.NewRequest("GET", "/fizzbuzz", nil)
	if newHttpRequestErr != nil {
		t.Fatal("Unexpected error:", newHttpRequestErr)
	}

	// Call the map handler on the HTTP request
	responseRecorder := httptest.NewRecorder()
	mapHandler(responseRecorder, httpRequest)
	response := responseRecorder.Result()

	// Examine response code, ensure it is expected
	code := response.StatusCode
	if code != http.StatusOK {
		t.Fatalf("Expected %d error code, got %d", http.StatusOK, code)
	}

	responseBodyBytes, responseBodyErr := ioutil.ReadAll(response.Body)
	if responseBodyErr != nil {
		t.Fatal("Unexpected error", responseBodyErr)
	}
	defer response.Body.Close()
	responseBody := strings.Trim(string(responseBodyBytes), "\r\n")
	if responseBody != HelloResponseBody {
		t.Fatalf("Wrong response body on default response.  Expected: \"%s\", found: \"%s\"", HelloResponseBody, responseBody)
	}
}

func TestMapHandlerNonDefault(t *testing.T) {
	mux := defaultMux()
	mapHandler := MapHandler(map[string]string{
		"/fizzbuzz": "https://localhost:420",
	}, mux)

	// Create HTTP request
	httpRequest, newHttpRequestErr := http.NewRequest("GET", "/fizzbuzz", nil)
	if newHttpRequestErr != nil {
		t.Fatal("Unexpected error:", newHttpRequestErr)
	}

	// Call map handler on the HTTP request
	responseRecorder := httptest.NewRecorder()
	mapHandler(responseRecorder, httpRequest)
	response := responseRecorder.Result()

	// Examine response code, ensure it is expected
	responseCode := response.StatusCode
	if responseCode != http.StatusSeeOther {
		t.Fatalf("Expected %d error code, got %d", http.StatusSeeOther, responseCode)
	}

	// Examine response header and location is as expected
	responseHeader := response.Header
	responseLocation := responseHeader["Location"][0]
	if responseLocation != "https://localhost:420" {
		t.Fatalf("Expected Location https://localhost:420, got %s", responseLocation)
	}
}

func TestMapHandlerWrongMethod(t *testing.T) {
	mux := defaultMux()
	mapHandler := MapHandler(map[string]string{}, mux)

	// Come up with an HTTP request
	httpRequest, newHttpRequestErr := http.NewRequest("POST", "/fizzbuzz", nil)
	if newHttpRequestErr != nil {
		t.Fatal("Unexpected error:", newHttpRequestErr)
	}

	// Call map handler on the HTTP request
	responseRecorder := httptest.NewRecorder()
	mapHandler(responseRecorder, httpRequest)
	response := responseRecorder.Result()

	// Examine response code, ensure it is expected
	code := response.StatusCode
	if code != http.StatusOK {
		t.Fatalf("Expected %d error code, got %d", http.StatusOK, code)
	}

	// Response should be the fallback
	responseBodyBytes, responseBodyErr := ioutil.ReadAll(response.Body)
	if responseBodyErr != nil {
		t.Fatal("Unexpected error", responseBodyErr)
	}
	defer response.Body.Close()
	responseBody := strings.Trim(string(responseBodyBytes), "\r\n")
	if responseBody != HelloResponseBody {
		t.Fatalf("Wrong response body on default response.  Expected: \"%s\", found: \"%s\"", HelloResponseBody, responseBody)
	}
}