package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/herrjulz/jefe/pkg/forms"

	"github.com/google/go-github/github"
	"github.com/herrjulz/jefe/pkg/models"
	"golang.org/x/oauth2"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	envs, err := app.Envs.List()
	if err != nil {
		log.Fatal(err)
		return
	}

	app.render(w, r, "home.page.tmpl", templateData{
		Environments:  envs,
		GitHubAuth:    app.AuthURL,
		User:          app.Session.GetString(r, "user"),
		Avatar:        app.Session.GetString(r, "avatar"),
		Authenticated: app.isAuthenticated(r),
	})
}

func (app *application) welcome(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "welcome.page.tmpl", templateData{
		GitHubAuth: app.AuthURL,
	})
}

func (app *application) claim(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	user := app.Session.GetString(r, "user")
	err = app.Envs.Claim(id, user)
	if err != nil {
		if errors.Is(err, models.ErrClaimed) {
			app.Session.Put(r, "flash", "You already claimed an environment")
		} else {
			app.Session.Put(r, "flash", err.Error())
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) unclaim(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	user := app.Session.GetString(r, "user")
	err = app.Envs.Unclaim(id, user)
	if err != nil {
		if errors.Is(err, models.ErrUnclaimed) {
			app.Session.Put(r, "flash", "You can just unclaim clusters you claimed")
		} else {
			app.Session.Put(r, "flash", err.Error())
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) unclaimAll(w http.ResponseWriter, r *http.Request) {
	err = app.Envs.UnclaimAll()
	if err != nil {
		app.Session.Put(r, "flash", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) createForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", templateData{
		GitHubAuth:    app.AuthURL,
		User:          app.Session.GetString(r, "user"),
		Authenticated: app.isAuthenticated(r),
		Avatar:        app.Session.GetString(r, "avatar"),
		Form:          forms.New(nil),
	})
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("envName", "envImage", "envAbout")
	form.MaxLength("envAbout", 230)

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", templateData{
			GitHubAuth:    app.AuthURL,
			User:          app.Session.GetString(r, "user"),
			Authenticated: app.isAuthenticated(r),
			Avatar:        app.Session.GetString(r, "avatar"),
			Form:          form,
		})
		return
	}

	err = app.Envs.Create(
		r.PostForm.Get("envName"),
		r.PostForm.Get("envImage"),
		r.PostForm.Get("envAbout"),
	)
	if err != nil {
		fmt.Println("err:", err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	app.Session.Put(r, "flash", fmt.Sprintf("Successfully created environment %s", r.PostForm.Get("envName")))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) deleteEnv(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		fmt.Println("err:", err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	err = app.Envs.Delete(id)
	if err != nil {
		fmt.Println("err:", err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	app.Session.Put(r, "flash", "Environment successfully deleted")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) oauthRedirect(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	code := r.FormValue("code")

	token, err := app.OAuthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("there was an issue getting your token. Err:", err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if !token.Valid() {
		fmt.Println(w, "retreived invalid token")
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	client := github.NewClient(app.OAuthConf.Client(oauth2.NoContext, token))

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		fmt.Println("error getting name")
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	orgs, _, err := client.Organizations.List(context.Background(), "", nil)
	if err != nil {
		fmt.Println("error getting name", err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if !app.userIsOrgMember(orgs) {
		app.Session.Put(r, "flash", fmt.Sprintf("You are not a member of the %s organization.", app.GHOAuthOrg))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r, "user", *user.Login)
	app.Session.Put(r, "avatar", *user.AvatarURL)
	app.Session.Put(r, "token", token.AccessToken)
	app.Session.Put(r, "flash", fmt.Sprintf("Hello %s! Welcome to the new claiming experience...", *user.Login))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Destroy(r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
