package config

import "html/template"

// AppConfig hold the application config
type AppConfig struct {
	TemplateCache map[string]*template.Template
}
