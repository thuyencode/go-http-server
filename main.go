package main

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

const PORT = 3000

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", HandleRootEndpoint)
	mux.HandleFunc("/goodbye", HandleGoodbyeEndpoint)
	mux.HandleFunc("/hello", HandleHelloEndpoint)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
	slog.Info("Listening to", "port", PORT)
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
	params := req.URL.Query()
	nameList, ok := params["name"]

	if !ok {
		http.Error(res, "You need to provide a search query named \"name\"", http.StatusBadRequest)
		return
	}

	name := nameList[0]
	var output bytes.Buffer

	output.WriteString("Hello, ")
	output.WriteString(name)
	output.WriteString("!")

	_, err := res.Write(output.Bytes())

	if err != nil {
		slog.Error("error writing response body:", "err", err)
		return
	}
}
