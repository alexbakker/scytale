package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	templates = map[string]*template.Template{}
)

func loadTemplates() error {
	baseLayout := "./templates/base.html"
	pages, err := filepath.Glob("./templates/pages/*.html")
	if err != nil {
		return err
	}

	for _, page := range pages {
		templates[filepath.Base(page)] =
			template.Must(template.ParseFiles(page, baseLayout))
	}

	return nil
}

func renderTemplate(w http.ResponseWriter, name string) error {
	tmpl, exists := templates[name]
	if !exists {
		return fmt.Errorf("template %s does not exist", name)
	}

	w.Header().Set("Content-Type", "text/html")
	return tmpl.ExecuteTemplate(w, "base", nil)
}
