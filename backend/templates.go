package main

import (
	"embed"
	"html/template"
	"log"
	"strings"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))

// render executes a template from templates/ without its trailing newline.
func render(name string, data any) string {
	var b strings.Builder
	if err := templates.ExecuteTemplate(&b, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
	}
	return strings.TrimSuffix(b.String(), "\n")
}
