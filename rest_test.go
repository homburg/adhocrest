package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testMethod(method string, callback func (w *httptest.ResponseRecorder)) {
	req, err := http.NewRequest(method, "http://example.com/foo", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	index(w, req)

	callback(w)
}

func TestInvalidMethod(t *testing.T) {
	testMethod("fisk", func (w *httptest.ResponseRecorder) {
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d\n", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestGet(t *testing.T) {
	testMethod("GET", func (w *httptest.ResponseRecorder) {
		if w.Body.String() != "GET\n" {
			t.Errorf("Expected body \"GET\", got\n-----\n%s\n----\n", string(w.Body.String()))
		}
	})
}
