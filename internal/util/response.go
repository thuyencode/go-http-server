package util

import (
	"bytes"
	"log/slog"
	"net/http"
)

func WriteResponseBody(res http.ResponseWriter, strings ...string) {
	var body bytes.Buffer

	for _, string := range strings {
		body.WriteString(string)
	}

	_, err := res.Write(body.Bytes())

	if err != nil {
		slog.Error("error writing response body:", "err", err)
		http.Error(res, "error writing response body", http.StatusInternalServerError)
		return
	}
}
