package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app application) newRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /connect", app.connect)

	standard := alice.New(app.recoverPanic)
	return standard.Then(mux)
}
