package router

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

type Router struct {
	funcs map[string]http.HandlerFunc
	// middleware []MiddlewareFunc
}

func New() *Router {
	return &Router{
		funcs: make(map[string]http.HandlerFunc, 0),
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	log.Info().Msg(
		fmt.Sprintf("{\"method\": \"%s\", \"path\": \"%s\"}", method, path))
	handler, ok := rt.funcs[method+path]
	if !ok {
		rt.handleNotFound(w, r)
		return
	}
	handler(w, r)
}

func (rt *Router) Get(path string, handler http.HandlerFunc) {
	internalPath := http.MethodGet + path
	rt.funcs[internalPath] = handler
}

func (rt *Router) Post(path string, handler http.HandlerFunc) {
	internalPath := http.MethodPost + path
	rt.funcs[internalPath] = handler
}

func (rt *Router) Put(path string, handler http.HandlerFunc) {
	internalPath := http.MethodPut + path
	rt.funcs[internalPath] = handler
}

func (rt *Router) Delete(path string, handler http.HandlerFunc) {
	internalPath := http.MethodDelete + path
	rt.funcs[internalPath] = handler
}

func (rt *Router) handleNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("route not found"))
}

func (rt *Router) Start(addr string) error {
	log.Info().Msg(fmt.Sprintf("starting server at %s", addr))
	if err := http.ListenAndServe(addr, rt); err != nil {
		return err
	}
	return nil
}
