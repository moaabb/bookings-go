package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/moaabb/bookings-go/internal/config"
	"github.com/moaabb/bookings-go/internal/handlers"
	"github.com/moaabb/bookings-go/internal/models"
	"github.com/moaabb/bookings-go/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager

const portNumber = ":8080"

func main() {

	// running the server config
	run()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	_ = http.ListenAndServe(portNumber, routes())

}

func run() {
	// Register data type for sessions
	gob.Register(models.Reservation{})

	// handling the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// app config
	app.InProduction = false
	app.Session = session

	// Creating the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot load template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

}
