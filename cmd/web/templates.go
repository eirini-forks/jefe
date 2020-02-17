package main

import (
	"text/template"
	"time"

	"github.com/herrjulz/jefe/pkg/forms"
	"github.com/herrjulz/jefe/pkg/models"
)

type templateData struct {
	Environments  []*models.Environment
	GitHubAuth    string
	User          string
	Authenticated bool
	Avatar        string
	Flash         string
	Form          *forms.Form
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
