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
	"github.com/moaabb/bookings-go/internal/driver"
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
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	_ = http.ListenAndServe(portNumber, routes(&app))

}

func run() (*driver.DB, error) {
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

	// connect to db
	log.Println("Connecting to Database...")
	db, err := driver.ConnectSQL("postgres://moab:example@localhost:8000/bookings")
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot connect to database! Dying...\tERROR: %s", err))
	}

	log.Println("Connected to Database!")

	// Creating the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot load template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	helpers.NewHelpers(&app)
	render.NewRenderer(&app)
	handlers.NewHandlers(repo)

	return db, nil

}
