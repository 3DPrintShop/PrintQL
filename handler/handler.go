// Package handler defines the HTTP handlers for this GraphQL API.
package handler

import (
	"bytes"
	"net/http"
)

func respond(w http.ResponseWriter, body []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	header := w.Header()
	header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	header.Set("Access-Control-Allow-Credentials", "true")
	header.Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.WriteHeader(200)
	_, _ = w.Write(body)
}

func isSupported(method string) bool {
	return method == "POST" || method == "GET" || method == "OPTIONS"
}

func errorJSON(msg string) []byte {
	buf := bytes.Buffer{}
	return buf.Bytes()
}
