package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/moaabb/go-web-dev/pkg/config"
	"github.com/moaabb/go-web-dev/pkg/handlers"
	"github.com/moaabb/go-web-dev/pkg/render"
)

var app config.AppConfig

const portNumber = ":8080"

func main() {
	// Hello world, the web server

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot load template cache", err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	_ = http.ListenAndServe(portNumber, nil)
}
