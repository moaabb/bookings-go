package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/moaabb/bookings-go/internal/config"
	"github.com/moaabb/bookings-go/internal/helpers"
	"github.com/moaabb/bookings-go/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func TestMain(m *testing.M) {
	// Register data type for sessions
	gob.Register(models.Reservation{})

	// when in production change this to true
	testApp.InProduction = false
	testApp.UseCache = false
	// handling the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	helpers.NewHelpers(&testApp)

	// app config
	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (mw *myWriter) Header() http.Header {
	var h http.Header

	return h
}

func (mw *myWriter) WriteHeader(StatusCode int) {

}

func (mw *myWriter) Write(b []byte) (int, error) {
	length := len(b)

	return length, nil
}
