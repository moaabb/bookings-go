package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/moaabb/go-web-dev/pkg/handlers"
)

func routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(traceUrl)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
