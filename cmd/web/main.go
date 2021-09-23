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

	tc, err := render.RenderTemplateCache()
	if err != nil {
		log.Fatal("Cannot load template cache", err)
	}

	app.TemplateCache = tc

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	_ = http.ListenAndServe(portNumber, nil)
}
