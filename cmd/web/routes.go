package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moaabb/bookings-go/internal/config"
	"github.com/moaabb/bookings-go/internal/handlers"
)

// routes manages the website routes
func routes(a *config.AppConfig) http.Handler {
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

	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/user/login", handlers.Repo.GetLogin)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	mux.Get("/rooms/generals", handlers.Repo.Generals)
	mux.Get("/rooms/majors", handlers.Repo.Majors)

	mux.Route("/admin", func(r chi.Router) {
		// r.Use(Auth)

		r.Get("/dashboard", handlers.Repo.AdminDashboard)

		r.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		r.Get("/reservations-new", handlers.Repo.AdminNewReservations)

		r.Get("/reservation-calendar", handlers.Repo.AdminReservationsCalendar)
		r.Post("/reservation-calendar", handlers.Repo.AdminPostReservationsCalendar)

		r.Get("/process-reservation/{{src}}/{{id}}", handlers.Repo.AdminProcessReservation)
		r.Get("/delete-reservation/{{src}}/{{id}}", handlers.Repo.AdminDeleteReservation)

		r.Get("/reservation/{src}/{id}", handlers.Repo.AdminShowReservation)
		r.Post("/reservation/{src}/{id}", handlers.Repo.AdminPostShowReservation)

	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.NotFound(handlers.Repo.NotFound)

	return mux
}
