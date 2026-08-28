package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnit_HandleRootEndpoint(t *testing.T) {
	res := httptest.NewRecorder()

	HandleRootEndpoint(res, nil)

	expectedCode := http.StatusOK
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d", expectedCode, res.Code)
	}

	expectedBody := []byte("Welcome!")
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf("bad response body, expected \"%s\", got \"%s\"", expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleGoodbyeEndpoint(t *testing.T) {
	res := httptest.NewRecorder()

	HandleGoodbyeEndpoint(res, nil)

	expectedCode := http.StatusOK
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d", expectedCode, res.Code)
	}

	expectedBody := []byte("Goodbye!")
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf("bad response body, expected \"%s\", got \"%s\"", expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleHelloEndpoint(t *testing.T) {
	name := "Gopher"
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/hello?name=%s", name), nil)

	HandleHelloEndpoint(res, req)

	expectedCode := http.StatusOK
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d", expectedCode, res.Code)
	}

	expectedBody := fmt.Appendf(nil, "Hello, %s!", name)
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf("bad response body, expected \"%s\", got \"%s\"", expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleHelloEndpoint_NoQuery(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)

	HandleHelloEndpoint(res, req)

	expectedCode := http.StatusBadRequest
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d", expectedCode, res.Code)
	}

	expectedBody := []byte("You need to provide a search query named \"name\"\n")
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf("bad response body, expected \"%s\", got \"%s\"", expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleHeaderEndpoint(t *testing.T) {
	name := "Gopher"
	req := httptest.NewRequest(http.MethodGet, "/header", nil)
	res := httptest.NewRecorder()

	req.Header.Set("name", name)
	HandleHeaderEndpoint(res, req)

	expectedCode := http.StatusOK
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d", expectedCode, res.Code)
	}

	expectedBody := fmt.Appendf(nil, "Hello, %s!", name)
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf("bad response body, expected \"%s\", got \"%s\"", expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleHeaderEndpoint_NoHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/header", nil)
	res := httptest.NewRecorder()

	HandleHeaderEndpoint(res, req)

	expectedCode := http.StatusBadRequest
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d", expectedCode, res.Code)
	}

	expectedBody := []byte("You must set a value for \"name\" at the request header\n")
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf("bad response body, expected \"%s\", got \"%s\"", expectedBody, res.Body.Bytes())
	}
}
