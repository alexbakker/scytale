package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	templateMap = loadTemplates()
	assetMap    = GetAssets()
)

func loadTemplates() (templates map[string]*template.Template) {
	templates = map[string]*template.Template{}
	baseLayout := "templates/base.html"
	pages := matchAssetByPrefix("templates/pages/")

	for _, page := range pages {
		//parse the child layout
		tmplName := filepath.Base(page)
		tmpl := template.Must(template.New(tmplName).Parse(string(assetMap[page])))

		//and finally also parse the base layout
		templates[tmplName] = template.Must(tmpl.Parse(string(assetMap[baseLayout])))
	}
	return
}

func matchAssetByPrefix(prefix string) (matches []string) {
	for key := range assetMap {
		if strings.HasPrefix(key, prefix) {
			matches = append(matches, key)
		}
	}
	return
}

func renderTemplate(w http.ResponseWriter, name string) error {
	tmpl, exists := templateMap[name]
	if !exists {
		return fmt.Errorf("template %s does not exist", name)
	}

	w.Header().Set("Content-Type", "text/html")
	return tmpl.ExecuteTemplate(w, "base", nil)
}
