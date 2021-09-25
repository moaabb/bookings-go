package render

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/moaabb/bookings-go/pkg/config"
	"github.com/moaabb/bookings-go/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// renderTemplate renders a html page
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		t, err := CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
		tc = t
	}

	parsedTemplate, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not load template")
	}
	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err := parsedTemplate.Execute(w, td)
	if err != nil {
		log.Println("error parsing template", err)
		return
	}
}

// CreateTemplateCache parses all the templates available and caches then
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		log.Println(err)
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
