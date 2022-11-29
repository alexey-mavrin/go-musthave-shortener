package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newServer(sh *store) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Get("/{key}", sh.fetchHandler)
	r.Post("/", sh.storeHandler)
	r.Post("/api/shorten", sh.storeJSONHandler)

	return r
}

func Run() error {
	sh := newStore()
	r := newServer(&sh)
	return http.ListenAndServe(":8080", r)
}
