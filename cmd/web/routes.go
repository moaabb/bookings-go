package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moaabb/bookings-go/internal/handlers"
)

// routes manages the website routes
func routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(traceUrl)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/rooms/generals", handlers.Repo.Generals)
	mux.Get("/rooms/majors", handlers.Repo.Majors)

	fileServer := http.FileServer(http.Dir("./static/"))
	fmt.Println(fileServer)
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
