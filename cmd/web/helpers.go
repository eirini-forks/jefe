package main

import (
	"bytes"
	"net/http"
	"path"
	"path/filepath"
	"text/template"

	"github.com/google/go-github/github"
)

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td templateData) {
	file := filepath.Join("./ui/html/", name)
	ts, err := template.New(path.Base(file)).Funcs(functions).ParseFiles(file, "./ui/html/base.layout.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err = ts.Execute(buf, app.addDefaultData(&td, r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.Flash = app.Session.PopString(r, "flash")
	return td
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r, "user")
}

func (app *application) userIsOrgMember(orgs []*github.Organization) bool {
	for _, org := range orgs {
		if *org.Login == app.GHOAuthOrg {
			return true
		}
	}
	return false
}
