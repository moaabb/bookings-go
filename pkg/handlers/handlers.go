package handlers

import (
	"net/http"

	"github.com/moaabb/go-web-dev/pkg/render"
)

func Home(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "index.page.tmpl")

}

func About(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl")
}
