package render

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

// renderTemplate renders a html page
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	tc, err := renderTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	parsedTemplate := tc[tmpl]
	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err = parsedTemplate.Execute(w, nil)
	if err != nil {
		log.Println("error parsing template", err)
		return
	}
}

func RenderTemplateCache() (map[string]*template.Template, error) {
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
