package shortener

import (
	"fmt"
	"github.com/gabbifish/urlshort/shortener/mocks"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	yaml := `
    - path: /urlshort
    url: https://github.com/gophercises/urlshort
    -
    url: https://github.com/gophercises/urlshort/tree/solution`
	_, err := NewYAMLSlugToURLFromBytes([]byte(yaml))

	if err == nil {
		t.Fatal("Incorrectly formed YAML file accepted.")
	}
}

func TestMapHandlerFallback(t *testing.T) {
	mux := defaultMux()

	slugToURL := new(mocks.SlugToURL)
	slugToURL.On("URLExists", "/fizzbuzz").Return(false, nil)
	slugToURL.On("URL", "/fizzbuzz").Return("", nil)
	httpHandler := NewHandlerFromSlugToURLClient(slugToURL, mux)

	// Come up with an HTTP request
	httpRequest, newHttpRequestErr := http.NewRequest("GET", "/fizzbuzz", nil)
	if newHttpRequestErr != nil {
		t.Fatal("Unexpected error:", newHttpRequestErr)
	}

	// Call the map handler on the HTTP request
	responseRecorder := httptest.NewRecorder()
	httpHandler.ServeHTTP(responseRecorder, httpRequest)
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
	slugToUrl := new(mocks.SlugToURL)
	slugToUrl.On("URLExists", "/fizzbuzz").Return(true, nil)
	slugToUrl.On("URL", "/fizzbuzz").Return("https://localhost:420", nil)

	//httpHandler := NewHandlerFromSlugToURLClient(slugToUrl, mux)
	//mapSlugToURL := NewMapSlugToURL(map[string]string{
	//	"/fizzbuzz": "https://localhost:420",
	//})
	mapHandler := NewHandlerFromSlugToURLClient(slugToUrl, mux)

	// Create HTTP request
	httpRequest, newHttpRequestErr := http.NewRequest("GET", "/fizzbuzz", nil)
	if newHttpRequestErr != nil {
		t.Fatal("Unexpected error:", newHttpRequestErr)
	}

	// Call map handler on the HTTP request
	responseRecorder := httptest.NewRecorder()
	mapHandler.ServeHTTP(responseRecorder, httpRequest)
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
	slugToURL := new(mocks.SlugToURL)
	slugToURL.On("URLExists", "/fizzbuzz").Return(false, nil)
	slugToURL.On("URL", "/fizzbuzz").Return("", nil)

	httpHandler := NewHandlerFromSlugToURLClient(slugToURL, mux)

	// Come up with an HTTP request
	httpRequest, newHttpRequestErr := http.NewRequest("POST", "/fizzbuzz", nil)
	if newHttpRequestErr != nil {
		t.Fatal("Unexpected error:", newHttpRequestErr)
	}

	// Call map handler on the HTTP request
	responseRecorder := httptest.NewRecorder()
	httpHandler.ServeHTTP(responseRecorder, httpRequest)
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

func TestHandlerWithNoMatchAndNoFallback(t *testing.T) {
	// Create a mocked SlugToURL - it will say that no URL exists in its mappings
	slugToURL := new(mocks.SlugToURL)
	slugToURL.On("URLExists", mock.Anything).Return(false, nil)
	slugToURL.On("URL", mock.Anything).Return("", nil)

	// Create a handler using our mocked SlugToURL, and with a nil backup http.Handler
	httpHandler := NewHandlerFromSlugToURLClient(slugToURL, nil)

	// Simulate a GET request to the slug "/fizzbuzz"
	responseRecorder := httptest.NewRecorder()
	httpRequest, newHttpRequestErr := http.NewRequest("GET", "/fizzbuzz", nil)
	if newHttpRequestErr != nil {
		t.Fatal("Unexpected error:", newHttpRequestErr)
	}
	httpHandler.ServeHTTP(responseRecorder, httpRequest)

	response := responseRecorder.Result()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status code %d but found %d", http.StatusNotFound, response.StatusCode)
	}
}

func TestHandlerWithMatchAndNoFallback(t *testing.T) {
	// Create a mocked SlugToURL - it will say that no URL exists in its mappings
	slugToURL := new(mocks.SlugToURL)
	slugToURL.On("URLExists", mock.Anything).Return(true, nil)
	slugToURL.On("URL", mock.Anything).Return("https://google.com", nil)

	// Create a handler using our mocked SlugToURL, and with a nil backup http.Handler
	httpHandler := NewHandlerFromSlugToURLClient(slugToURL, nil)

	// Simulate a GET request to the slug "/fizzbuzz"
	responseRecorder := httptest.NewRecorder()
	httpRequest, newHttpRequestErr := http.NewRequest("GET", "/fizzbuzz", nil)
	if newHttpRequestErr != nil {
		t.Fatal("Unexpected error:", newHttpRequestErr)
	}
	httpHandler.ServeHTTP(responseRecorder, httpRequest)

	response := responseRecorder.Result()
	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("Expected status code %d but found %d", http.StatusSeeOther, response.StatusCode)
	}
}
