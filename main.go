package main

import (
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

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
	slog.Info("Listening to", "port", PORT)
}

func HandleRootEndpoint(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello World!"))

	if err != nil {
		slog.Error("error writing response:", "err", err)
		return
	}
}

func HandleGoodbyeEndpoint(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Goodbye!"))

	if err != nil {
		slog.Error("error writing response:", "err", err)
		return
	}
}
