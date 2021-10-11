package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/moaabb/bookings-go/internal/config"
	"github.com/moaabb/bookings-go/internal/driver"
	"github.com/moaabb/bookings-go/internal/forms"
	"github.com/moaabb/bookings-go/internal/helpers"
	"github.com/moaabb/bookings-go/internal/models"
	"github.com/moaabb/bookings-go/internal/render"
	"github.com/moaabb/bookings-go/internal/repository"
	"github.com/moaabb/bookings-go/internal/repository/dbrepo"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

// NewRepo creates a respostiory and parse the data to the repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home renders the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "index.page.tmpl", &models.TemplateData{})
}

// Aboute renders the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// SearchAvailability renders the SearchAvailability page
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// Generals renders the Generals room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the Majors room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// MakeReservation renders the MakeReservation page
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, r.Form.Get("start_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, r.Form.Get("end_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room_id, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		RoomID:    room_id,
		StartDate: startDate,
		EndDate:   endDate,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	reservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        room_id,
		ReservationID: reservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Tipo    string `json:"type"`
}

// PostAvailability handles room availability form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	fmt.Println("chega")

	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

// AvailabilityJSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		OK:      true,
		Message: "Unavailable!",
		Tipo:    "error",
	}

	out, err := json.Marshal(resp)
	if err != nil {
		helpers.ServerError(w, err)
	}

	fmt.Println(r.Form.Get("start"), r.Form.Get("end"))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Could not fetch reservation data from session")
		m.App.Session.Put(r.Context(), "Error", "Could not fetch reservation data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data["reservation"] = reservation
	m.App.Session.Remove(r.Context(), "reservation")

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
