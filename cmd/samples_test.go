package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handleHome(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	handleHome(w, r)

	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("got %d but expected %d\n", result.StatusCode, http.StatusOK)
	}
}

func Test_handlePostThing(t *testing.T) {
	expectedResponse := `{"message":"1 things now"}`
	testBody := strings.NewReader(`{"id":1,"name":"ken"}`)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/thing", testBody)

	handlePostThing(w, r)
	result := w.Result()
	if result.StatusCode != http.StatusCreated {
		t.Errorf("got %d but expected %d\n", result.StatusCode, http.StatusCreated)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)
	responseString := string(bts)
	if responseString != expectedResponse {
		t.Errorf("got %s but expected %s\n", responseString, expectedResponse)
	}
}

func Test_getThings(t *testing.T) {
	// expectedResponse := "[]"
	expectedHeader := "application/json"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/thing", nil)

	withJson(handleGetThings)(w, r)

	result := w.Result()
	header := result.Header["Content-Type"]
	if header[0] != expectedHeader {
		t.Errorf("got %s but expected %s\n", header[0], expectedHeader)
	}
	// defer result.Body.Close()
	// bts, _ := io.ReadAll(result.Body)
	// responseString := string(bts)
	// if responseString != expectedResponse {
	// 	t.Errorf("got %s but expected %s\n", responseString, expectedResponse)
	// }
}
