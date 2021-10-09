package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/moaabb/bookings-go/internal/config"
	"github.com/moaabb/bookings-go/internal/handlers"
	"github.com/moaabb/bookings-go/internal/helpers"
	"github.com/moaabb/bookings-go/internal/models"
	"github.com/moaabb/bookings-go/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

const portNumber = ":8080"

func main() {

	// running the server config
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	_ = http.ListenAndServe(portNumber, routes(&app))

}

func run() error {
	// Register data type for sessions
	gob.Register(models.Reservation{})

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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
	helpers.NewHelpers(&app)
	render.NewTemplates(&app)
	handlers.NewHandlers(repo)

	return nil

}
