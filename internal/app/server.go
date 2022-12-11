package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newServer(c Config) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Get("/{key}", c.fetchHandler)
	r.Post("/", c.storeHandler)
	r.Post("/api/shorten", c.storeJSONHandler)

	return r
}

func Run() error {
	c := Config{
		ServerAddress: ":8080",
		BaseURL:       "http://localhost:8080/",
	}
	return c.Run()
}

func (c Config) Run() error {
	c.sh = newStore()
	r := newServer(c)
	return http.ListenAndServe(c.ServerAddress, r)
}
