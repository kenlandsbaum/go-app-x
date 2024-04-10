package main

import (
	"encoding/json"
	"fmt"
	"go-app-x/internal/router"
	"log"
	"net/http"
)

var (
	things = []Thing{}
)

type Thing struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func App() {
	r := router.New()
	r.Get("/", withJson(handleHome))
	r.Post("/thing", withJson(handlePostThing))
	r.Get("/thing", withJson(handleGetThings))
	if err := r.Start(":8999"); err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\":\"welcome to the api\"}"))
}

func handlePostThing(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	var thing Thing
	decoder := json.NewDecoder(body)
	decoder.Decode(&thing)
	things = append(things, thing)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("{\"message\":\"%d things now\"}", len(things))))
}

func handleGetThings(w http.ResponseWriter, r *http.Request) {
	bts, err := json.Marshal(things)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"message\":\"error getting things\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bts)
}

func withJson(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fn(w, r)
	}
}
