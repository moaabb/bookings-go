package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
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

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
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

// GetLogin renders the login page
func (m *Repository) GetLogin(w http.ResponseWriter, r *http.Request) {
	if m.App.Session.Exists(r.Context(), "user_id") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostLogin Handles the login Form post
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println("Error parsing login form")
		m.App.Session.Put(r.Context(), "error", "Error parsing form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials!")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in Successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// MakeReservation renders the MakeReservation page and displays form
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Could not fetch reservation data from session")
		m.App.Session.Put(r.Context(), "error", "Could not fetch reservation data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Could not fetch reservation data from session")
		m.App.Session.Put(r.Context(), "error", "Could not fetch reservation data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = res
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	reservationID, err := m.DB.InsertReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		RoomID:        res.RoomID,
		ReservationID: reservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"msg"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    int    `json:"room_id"`
}

// PostAvailability handles room availability form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	layout := "2006-01-02"

	start, err := time.Parse(layout, r.Form.Get("start"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	end, err := time.Parse(layout, r.Form.Get("end"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(start, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find Room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	for _, v := range rooms {
		fmt.Println("ROOM", v.ID, v.RoomName)

	}

	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No Availability!")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: start,
		EndDate:   end,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

// AvailabilityJSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room_id, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, room_id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find Room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		RoomID:    room_id,
		StartDate: sd,
		EndDate:   ed,
	}

	out, err := json.Marshal(&resp)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Could not fetch reservation data from session")
		m.App.Session.Put(r.Context(), "error", "Could not fetch reservation data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := res.StartDate.Format("Mon, Jan 02, 2006")
	ed := res.EndDate.Format("Mon, Jan 02, 2006")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data["reservation"] = res
	fmt.Println(res.Room.RoomName)
	m.App.Session.Remove(r.Context(), "reservation")

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Could not fetch reservation data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomID

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find Room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) NotFound(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "404.page.tmpl", &models.TemplateData{})
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation
	res.StartDate = startDate
	res.EndDate = endDate
	res.RoomID = roomID
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// Logout Logs the User out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	if !m.App.Session.Exists(r.Context(), "user_id") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	m.App.Session.Put(r.Context(), "flash", "Successfully logged out!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AdminDashboard Renders the Admin Tool
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminAllReservations Lists All Reservations in the Admin Tool
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	res, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = res

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminNewReservations Lists All New Reservations in the Admin Tool
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	res, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = res

	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// assume that there's no month/year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))

		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)

	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear

	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("January")
	stringMap["this_month_year"] = now.Format("2006")

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)

	intMap["days_of_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		IntMap:    intMap,
		Data:      data,
	})
}

// AdminShowReservation Renders the Edit Rservation page for Admins
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	resID, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	reservation, err := m.DB.GetReservationByID(resID)
	if err != nil {
		m.App.ErrorLog.Println("Can't find Reservation")
		m.App.Session.Put(r.Context(), "error", "Can't find Reservation")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservation-%s", src), http.StatusSeeOther)
		return
	}

	stringMap := make(map[string]string)
	stringMap["src"] = src

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

// AdminPostShowReservation Handles the post of the Edit Screen
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	resID, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	res, err := m.DB.GetReservationByID(resID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Phone = r.Form.Get("phone")
	res.Email = r.Form.Get("email")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes Saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}

// AdminProcessReservation Mark Reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	src := exploded[3]

	err = m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		m.App.ErrorLog.Println("Error Marking reservation as processed")
		m.App.Session.Put(r.Context(), "error", "Error Marking reservation as processed")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation Marked as Processed")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminDeleteReservation Deletes a reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	src := exploded[3]

	err = m.DB.DeleteReservation(id)
	if err != nil {
		m.App.ErrorLog.Println("Error deleting reservation")
		m.App.Session.Put(r.Context(), "error", "Error deleting reservation")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation Deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}
