package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-http-server/internal/users"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"reflect"
	"testing"
)

func TestUnit_HandleRootEndpoint(t *testing.T) {
	res := httptest.NewRecorder()

	HandleRootEndpoint(res, nil)

	expectedCode := http.StatusOK
	if res.Code != expectedCode {
		t.Errorf(`bad response code, expected %d, got %d`, expectedCode, res.Code)
	}

	expectedBody := []byte("Welcome!")
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleGoodbyeEndpoint(t *testing.T) {
	res := httptest.NewRecorder()

	HandleGoodbyeEndpoint(res, nil)

	expectedCode := http.StatusOK
	if res.Code != expectedCode {
		t.Errorf(`bad response code, expected %d, got %d`, expectedCode, res.Code)
	}

	expectedBody := []byte("Goodbye!")
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleHelloEndpoint(t *testing.T) {
	name := "Gopher"
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/hello?name=%s", name), nil)

	HandleHelloEndpoint(res, req)

	expectedCode := http.StatusOK
	if res.Code != expectedCode {
		t.Errorf(`bad response code, expected %d, got %d`, expectedCode, res.Code)
	}

	expectedBody := fmt.Appendf(nil, "Hello, %s!", name)
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleHelloEndpoint_NoQuery(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)

	HandleHelloEndpoint(res, req)

	expectedCode := http.StatusBadRequest
	if res.Code != expectedCode {
		t.Errorf(`bad response code, expected %d, got %d`, expectedCode, res.Code)
	}

	expectedBody := []byte("You need to provide a search query named \"name\"\n")
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedBody, res.Body.Bytes())
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
		t.Errorf(`bad response code, expected %d, got %d`, expectedCode, res.Code)
	}

	expectedBody := fmt.Appendf(nil, "Hello, %s!", name)
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleHeaderEndpoint_NoHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/header", nil)
	res := httptest.NewRecorder()

	HandleHeaderEndpoint(res, req)

	expectedCode := http.StatusBadRequest
	if res.Code != expectedCode {
		t.Errorf(`bad response code, expected %d, got %d`, expectedCode, res.Code)
	}

	expectedBody := []byte("You must set a value for \"name\" at the request header\n")
	if !bytes.Equal(res.Body.Bytes(), expectedBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedBody, res.Body.Bytes())
	}
}

func TestUnit_HandleJSONEndpoint(t *testing.T) {
	reqBody := UserData{FirstName: "Gopher"}
	marshalledReqBody, err := json.Marshal(reqBody)

	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(marshalledReqBody))
	res := httptest.NewRecorder()

	HandleJSONEndpoint(res, req)

	expectedCode := http.StatusOK
	if res.Code != expectedCode {
		t.Errorf(`bad response code, expected %d, got %d`, expectedCode, res.Code)
	}

	expectedResponseBody := []byte("Hello, Gopher!")
	if !bytes.Equal(res.Body.Bytes(), expectedResponseBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedResponseBody, res.Body.Bytes())
	}
}

func TestUnit_HandleJSONEndpoint_EmptyRequestBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/json", nil)
	res := httptest.NewRecorder()

	HandleJSONEndpoint(res, req)

	expectedCode := http.StatusBadRequest
	if res.Code != expectedCode {
		t.Errorf(`bad response code, expected %d, got %d`, expectedCode, res.Code)
	}

	expectedResponseBody := []byte("error deserialising request body\n")
	if !bytes.Equal(res.Body.Bytes(), expectedResponseBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedResponseBody, res.Body.Bytes())
	}
}

func TestUnit_HandleJSONEndpoint_EmptyName(t *testing.T) {
	reqBody := UserData{FirstName: ""}
	marshalledReqBody, err := json.Marshal(reqBody)

	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(marshalledReqBody))
	res := httptest.NewRecorder()

	HandleJSONEndpoint(res, req)

	expectedCode := http.StatusBadRequest
	if res.Code != expectedCode {
		t.Errorf(`bad response code, expected %d, got %d`, expectedCode, res.Code)
	}

	expectedResponseBody := []byte("You must not leave the \"name\" field empty\n")
	if !bytes.Equal(res.Body.Bytes(), expectedResponseBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedResponseBody, res.Body.Bytes())
	}
}

func TestIntegration_HandleUserEndpointPOST(t *testing.T) {
	expectedUserData := UserData{FirstName: "Go", LastName: "Gopher", Email: "gopher@email.com"}
	marshalledReqBody, err := json.Marshal(expectedUserData)

	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(marshalledReqBody))
	res := httptest.NewRecorder()

	testManager := users.NewManager()
	testServer := Server{userManager: testManager}

	req.Header.Set("Content-Type", "application/json")
	testServer.HandleUserEndpointPOST(res, req)

	expectedCode := http.StatusCreated
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d,\nresponse body: \"%s\"", expectedCode, res.Code, res.Body.String())
	}

	resultUser, err := testManager.GetUserByName(expectedUserData.FirstName, expectedUserData.LastName)
	if err != nil {
		t.Fatalf("error getting user from manager: %v", err)
	}

	actualUserData := convertUserToUserData(resultUser)

	if !reflect.DeepEqual(actualUserData, &expectedUserData) {
		t.Fatalf("bad retrieved user,\nexpected: %v\n,got: %v", expectedUserData, actualUserData)
	}
}

func TestIntegration_HandleUserEndpointPOST_BadHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/user", nil)
	res := httptest.NewRecorder()

	testManager := users.NewManager()
	testServer := Server{userManager: testManager}

	testServer.HandleUserEndpointPOST(res, req)

	expectedCode := http.StatusUnsupportedMediaType
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d,\nresponse body: \"%s\"", expectedCode, res.Code, res.Body.String())
	}

	expectedResponseBody := []byte("unsupported Content-Type header: \"\"\n")
	if !bytes.Equal(res.Body.Bytes(), expectedResponseBody) {
		t.Errorf(`bad response body, expected "%s", got "%s"`, expectedResponseBody, res.Body.Bytes())
	}
}

func TestIntegration_HandleUserEndpointGET(t *testing.T) {
	firstName, lastName, email := "Go", "Gopher", "gopher@email.com"

	testManager := users.NewManager()
	testServer := Server{userManager: testManager}

	err := testManager.AddUser(firstName, lastName, email)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/user", nil)
	q := req.URL.Query()

	q.Add("firstName", firstName)
	q.Add("lastName", lastName)
	req.URL.RawQuery = q.Encode()

	testServer.HandleUserEndpointGET(res, req)

	expectedCode := http.StatusOK
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d,\nresponse body: \"%s\"", expectedCode, res.Code, res.Body.String())
	}

	decoder := json.NewDecoder(res.Body)
	decoder.DisallowUnknownFields()

	var result UserData
	err = decoder.Decode(&result)

	if err != nil {
		t.Fatalf("error decoding response body: %v", err)
	}

	expectedResult := UserData{firstName, lastName, email}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("bad result, expected: %v\n, got: %v", expectedResult, result)
	}

	expectedContentType := "application/json"
	actualContentType := res.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("bad result, expected: %q, got: %q", expectedContentType, actualContentType)
	}
}

func TestIntegration_HandleUserEndpointGET_EmptyQuery(t *testing.T) {
	firstName, lastName, email := "Go", "Gopher", "gopher@email.com"

	testManager := users.NewManager()
	testServer := Server{userManager: testManager}

	err := testManager.AddUser(firstName, lastName, email)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	testFirstName, testLastName := "", ""

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/user", nil)
	q := req.URL.Query()

	q.Add("firstName", testFirstName)
	q.Add("lastName", testLastName)
	req.URL.RawQuery = q.Encode()

	testServer.HandleUserEndpointGET(res, req)

	expectedCode := http.StatusBadRequest
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d,\nresponse body: \"%s\"", expectedCode, res.Code, res.Body.String())
	}

	expectedResponseBody := []byte("Netheir \"firstName\" or \"lastName\" search query can be empty\n")
	if !bytes.Equal(res.Body.Bytes(), expectedResponseBody) {
		t.Errorf("bad result, expected: %q, got: %q", expectedResponseBody, res.Body.String())
	}
}

func TestIntegration_HandleUserEndpointGET_NoUserFound(t *testing.T) {
	firstName, lastName, email := "Go", "Gopher", "gopher@email.com"

	testManager := users.NewManager()
	testServer := Server{userManager: testManager}

	err := testManager.AddUser(firstName, lastName, email)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	testFirstName, testLastName := "Rust", "Ferris"

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/user", nil)
	q := req.URL.Query()

	q.Add("firstName", testFirstName)
	q.Add("lastName", testLastName)
	req.URL.RawQuery = q.Encode()

	testServer.HandleUserEndpointGET(res, req)

	expectedCode := http.StatusNotFound
	if res.Code != expectedCode {
		t.Errorf("bad response code, expected %d, got %d,\nresponse body: \"%s\"", expectedCode, res.Code, res.Body.String())
	}

	expectedResponseBody := []byte("No user found\n")
	if !bytes.Equal(res.Body.Bytes(), expectedResponseBody) {
		t.Errorf("bad result, expected: %q, got: %q", expectedResponseBody, res.Body.String())
	}
}

func TestUnit_convertUserToUserData(t *testing.T) {
	firstName, lastName := "Go", "Gopher"
	email, err := mail.ParseAddress("gopher@email.com")

	if err != nil {
		t.Fatalf("error parsing email address: %v", err)
	}

	testUser := users.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     *email,
	}
	actualUserData := convertUserToUserData(&testUser)
	expectedUserData := &UserData{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email.Address,
	}

	if !reflect.DeepEqual(actualUserData, expectedUserData) {
		t.Errorf("bad conversion,\nexpected: %v,\ngot: %v", expectedUserData, actualUserData)
	}
}
