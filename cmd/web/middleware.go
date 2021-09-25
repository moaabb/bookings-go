package main

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

// traceUrl Logs the request to the server
func traceUrl(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.URL)
		next.ServeHTTP(w, req)
	})
}

// NoSurf Prevent CSRF Attacks
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves session data for current request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
