package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/berberapan/my-stuff/internal/models"
	"github.com/berberapan/my-stuff/ui"
)

type templateData struct {
	CurrentYear     int
	Item            models.Item
	Items           []models.Item
	Form            any
	IsAuthenticated bool
	CSRFToken       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.tmpl",
			"html/partials/navbar.tmpl",
			page,
		}

		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
