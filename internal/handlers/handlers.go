package handlers

import (
	"net/http"

	"github.com/moaabb/bookings-go/internal/config"
	"github.com/moaabb/bookings-go/internal/models"
	"github.com/moaabb/bookings-go/internal/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

// NewRepo creates a respostiory and parse the data to the repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home renders the home page
func (m *Repository) Home(w http.ResponseWriter, req *http.Request) {

	m.App.Session.Put(req.Context(), "remote_ip", req.RemoteAddr)

	render.RenderTemplate(w, "index.page.tmpl", &models.TemplateData{})

}

// Aboute renders the about page
func (m *Repository) About(w http.ResponseWriter, req *http.Request) {

	StringMap := make(map[string]string)

	StringMap["remote_ip"] = m.App.Session.GetString(req.Context(), "remote_ip")

	StringMap["test"] = "Olá de novo"
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: StringMap,
	})
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{})
}

// SearchAvailability renders the SearchAvailability page
func (m *Repository) SearchAvailability(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{})
}

// Generals renders the Generals room page
func (m *Repository) Generals(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the Majors room page
func (m *Repository) Majors(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{})
}

// MakeReservation renders the MakeReservation page
func (m *Repository) MakeReservation(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{})
}
