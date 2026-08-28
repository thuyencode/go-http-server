package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRootEndpoint(t *testing.T) {
	r := httptest.NewRecorder()
	HandleRootEndpoint(r, nil)

	expectedCode := http.StatusOK

	if r.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d", expectedCode, r.Code)
	}

	expectedBody := []byte("Hello World!")

	if !bytes.Equal(r.Body.Bytes(), expectedBody) {
		t.Errorf("bad response body, expected \"%s\", got \"%s\"", expectedBody, r.Body.Bytes())
	}
}

func TestHandleGoodbyeEndpoint(t *testing.T) {
	r := httptest.NewRecorder()
	HandleGoodbyeEndpoint(r, nil)

	expectedCode := http.StatusOK

	if r.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d", expectedCode, r.Code)
	}

	expectedBody := []byte("Goodbye!")

	if !bytes.Equal(r.Body.Bytes(), expectedBody) {
		t.Errorf("bad response body, expected \"%s\", got \"%s\"", expectedBody, r.Body.Bytes())
	}
}
