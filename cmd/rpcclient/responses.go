package main

import "net/http"

var (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

func handleFacePlant(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func handleSuccess(w http.ResponseWriter, bts []byte, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write(bts)
}

func withJsonContent(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(contentType, applicationJson)
		fn(w, r)
	}
}
