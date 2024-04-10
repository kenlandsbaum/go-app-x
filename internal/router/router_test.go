package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeReqEcho(method, body string) string {
	return fmt.Sprintf("method:%s,body:%s", method, body)
}

func testingGetHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	response := fmt.Sprintf("responding with path: %s, method: %s", path, method)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func testingPostHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	body := r.Body
	defer body.Close()

	bts, _ := io.ReadAll(body)
	response := makeReqEcho(method, string(bts))
	w.Write([]byte(response))
}

func TestNew(t *testing.T) {
	expectedBody := "responding with path: /test, method: GET"
	r := New()
	r.Get("/test", testingGetHandler)

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	result := rec.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("got %d but expected %d\n", rec.Result().StatusCode, http.StatusOK)
	}
	defer result.Body.Close()
	actualBody, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("expected nil error but got %s\n", err.Error())
	}

	if string(actualBody) != expectedBody {
		t.Errorf("got %s but expected %s\n", string(actualBody), expectedBody)
	}
}

func TestPost(t *testing.T) {
	expected := makeReqEcho("POST", "{}")
	r := New()
	r.Post("/test", testingPostHandler)

	postBody := strings.NewReader("{}")

	req, _ := http.NewRequest(http.MethodPost, "/test", postBody)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	result := rec.Result()
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if string(bts) != expected {
		t.Errorf("got %s but expected %s", string(bts), expected)
	}
}

func TestPut(t *testing.T) {
	expected := makeReqEcho("PUT", "{}")
	r := New()
	r.Put("/test", testingPostHandler)

	_, ok := r.funcs["PUT/test"]
	if !ok {
		t.Error("expected func for PUT/test")
	}

	postBody := strings.NewReader("{}")

	req, _ := http.NewRequest(http.MethodPut, "/test", postBody)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	result := rec.Result()
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if string(bts) != expected {
		t.Errorf("got %s but expected %s", string(bts), expected)
	}
}

func TestDelete(t *testing.T) {
	expected := makeReqEcho("DELETE", "{}")
	r := New()
	r.Delete("/test", testingPostHandler)

	postBody := strings.NewReader("{}")

	req, _ := http.NewRequest(http.MethodDelete, "/test", postBody)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	result := rec.Result()
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if string(bts) != expected {
		t.Errorf("got %s but expected %s", string(bts), expected)
	}
}

func TestNotFound(t *testing.T) {
	expected := "route not found"
	r := New()
	r.Put("/test", testingPostHandler)

	postBody := strings.NewReader("{}")

	req, _ := http.NewRequest(http.MethodGet, "/test", postBody)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	result := rec.Result()
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if string(bts) != expected {
		t.Errorf("got %s but expected %s", string(bts), expected)
	}
}
