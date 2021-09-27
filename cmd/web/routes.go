package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moaabb/bookings-go/pkg/handlers"
)

// routes manages the website routes
func routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(traceUrl)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	fileServer := http.FileServer(http.Dir("./static/"))
	fmt.Println(fileServer)
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
