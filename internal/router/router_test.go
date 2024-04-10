package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	expectedBody := "responding with path: /test, method: GET"
	r := New()
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		response := fmt.Sprintf("responding with path: %s, method: %s", path, method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

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
