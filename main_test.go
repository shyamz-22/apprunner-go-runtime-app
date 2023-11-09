package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type MockDB struct {
}

func (db *MockDB) SaveURL(longurl string) (string, error) {
	return "abcdef1234", nil
}

func (db *MockDB) GetLongURL(shortCode string) (string, error) {
	return "https://example.com", nil
}

func TestCreateShortCode(t *testing.T) {
	// Setup the router
	router := mux.NewRouter()
	server := NewUrlShortener(&MockDB{})
	router.HandleFunc("/app/", server.createShortCode).Methods(http.MethodPost)

	// When: create request to create a short code
	req, err := http.NewRequest(http.MethodPost, "/app/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Body = io.NopCloser(strings.NewReader("https://example.com"))

	// Make the request and record the response
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	// Then: validate status code
	if rw.Code != http.StatusCreated {
		body := string(rw.Body.Bytes())
		t.Errorf("expected status code %d, got %d\nerror: %s", http.StatusCreated, rw.Code, body)
		t.FailNow()
	}

	// Decode the response body
	var response CreateShortCodeResponse
	err = json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Then: Validate the response body
	if len(response.ShortCode) == 0 {
		t.Errorf("expected short code, got %s", response.ShortCode)
	}
}

func TestAccessURLWithShortCode(t *testing.T) {
	// Given: A shortened url
	shortcode := createShortenedUrl()

	// When: Access Url with short code
	router := mux.NewRouter()
	server := NewUrlShortener(&MockDB{})

	router.HandleFunc("/app/{shortcode}", server.accessURLWithShortCode).Methods(http.MethodGet)

	// Make a request to access the URL
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/app/%s", shortcode), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Make the request and record the response
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	// Then: Validate status code
	if rw.Code != http.StatusFound {
		body := string(rw.Body.Bytes())
		t.Errorf("expected status code %d, got %d\nerror:%s", http.StatusFound, rw.Code, body)
		t.FailNow()
	}

	// Then: Validate the Location header
	location := rw.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("expected Location header to be https://example.com, got %s", location)
	}
}

func createShortenedUrl() string {
	router := mux.NewRouter()
	server := NewUrlShortener(&MockDB{})

	router.HandleFunc("/app/", server.createShortCode).Methods(http.MethodPost)

	// Make a request to create a short code
	req, _ := http.NewRequest(http.MethodPost, "/app/", nil)
	req.Body = io.NopCloser(strings.NewReader("https://example.com"))

	// Make the request and record the response
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	var response CreateShortCodeResponse
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		_ = fmt.Errorf("failed to decode response body: %v", err)
		os.Exit(1)
	}

	return response.ShortCode

}
