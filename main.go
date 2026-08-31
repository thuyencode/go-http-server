package main

import (
	"encoding/json"
	"fmt"
	"go-http-server/internal/util"
	"io"
	"log"
	"log/slog"
	"net/http"
)

const PORT = 3000

type RequestBody struct {
	Name string
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", HandleRootEndpoint)
	mux.HandleFunc("/goodbye", HandleGoodbyeEndpoint)
	mux.HandleFunc("/hello", HandleHelloEndpoint)
	mux.HandleFunc("/param/{name}", HandleParamEndpoint)
	mux.HandleFunc("/header", HandleHeaderEndpoint)
	mux.HandleFunc("/json", HandleJSONEndpoint)

	slog.Info("Listening to", "port", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
}

func HandleRootEndpoint(res http.ResponseWriter, req *http.Request) {
	_, err := res.Write([]byte("Welcome!"))

	if err != nil {
		slog.Error("error writing response body:", "err", err)
		return
	}
}

func HandleGoodbyeEndpoint(res http.ResponseWriter, _ *http.Request) {
	_, err := res.Write([]byte("Goodbye!"))

	if err != nil {
		slog.Error("error writing response body:", "err", err)
		return
	}
}

func HandleHelloEndpoint(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	nameList, ok := query["name"]

	if !ok {
		http.Error(res, "You need to provide a search query named \"name\"", http.StatusBadRequest)
		return
	}

	name := nameList[0]

	util.WriteResponseBody(res, "Hello, ", name, "!")
}

func HandleParamEndpoint(res http.ResponseWriter, req *http.Request) {
	name := req.PathValue("name")
	util.WriteResponseBody(res, "Hello, ", name, "!")
}

func HandleHeaderEndpoint(res http.ResponseWriter, req *http.Request) {
	name := req.Header.Get("name")

	if name == "" {
		http.Error(res, "You must set a value for \"name\" at the request header", http.StatusBadRequest)
		return
	}

	util.WriteResponseBody(res, "Hello, ", name, "!")
}

func HandleJSONEndpoint(res http.ResponseWriter, req *http.Request) {
	bodyInBytes, err := io.ReadAll(req.Body)

	if err != nil {
		slog.Error("error reading request body:", "err", err)
		http.Error(res, "error reading request body", http.StatusInternalServerError)
		return
	}

	var unmarshalledRequestBody RequestBody
	err = json.Unmarshal(bodyInBytes, &unmarshalledRequestBody)

	if err != nil {
		slog.Error("error deserialising request body:", "err", err)
		http.Error(res, "error deserialising request body", http.StatusBadRequest)
		return
	}

	if unmarshalledRequestBody.Name == "" {
		http.Error(res, `You must not leave the "name" field empty`, http.StatusBadRequest)
		return
	}

	util.WriteResponseBody(res, "Hello, ", unmarshalledRequestBody.Name, "!")
}
