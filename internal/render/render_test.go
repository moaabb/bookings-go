package render

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/moaabb/bookings-go/internal/models"
)

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	// tc, err := CreateTemplateCache()
	// if err != nil {
	// 	t.Error(err)
	// }

	// app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "index.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to the browser", err)
	}

	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	testMsg := "test flash"

	r, err := getSession()
	if err != nil {
		t.Fatal(err)
	}

	session.Put(r.Context(), "Flash", testMsg)
	result := AddDefaultData(&td, r)

	if result.Flash != testMsg {
		t.Error(fmt.Sprintf("Flash msg, expected '%s', got '%s'", testMsg, result.Flash))
	}

}

func TestCreateTemplateCache(t *testing.T) {
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()

	ctx, err = session.Load(ctx, r.Header.Get("X-Session"))
	if err != nil {
		return nil, err
	}
	r = r.WithContext(ctx)

	return r, nil
}
