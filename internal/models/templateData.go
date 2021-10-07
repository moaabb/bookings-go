package models

import "github.com/moaabb/bookings-go/internal/forms"

// TemplateData is the type of data that will be parsed to the templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Error     string
	Warning   string
	Form      *forms.Form
}
