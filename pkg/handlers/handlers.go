package handlers

import (
	"net/http"

	"github.com/moaabb/go-web-dev/pkg/config"
	"github.com/moaabb/go-web-dev/pkg/models"
	"github.com/moaabb/go-web-dev/pkg/render"
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

// Home handles the home page
func (m *Repository) Home(w http.ResponseWriter, req *http.Request) {
	StringMap := make(map[string]string)

	StringMap["test"] = "Olá de novo"

	render.RenderTemplate(w, "index.page.tmpl", &models.TemplateData{
		StringMap: StringMap,
	})

}

// Aboute handles the about page
func (m *Repository) About(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl", nil)
}
