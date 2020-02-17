package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddlware := alice.New(app.secureHeaders, app.Session.Enable)
	dynamicMiddleware := alice.New(app.requireAuth)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.home)))
	mux.Get("/welcome", http.HandlerFunc(app.welcome))

	mux.Post("/envs/claim/:id", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.claim)))
	mux.Post("/envs/unclaim/:id", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.unclaim)))
	mux.Post("/envs/unclaim", http.HandlerFunc(app.unclaimAll))

	mux.Get("/envs/create", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.createForm)))
	mux.Post("/envs/create", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.create)))
	mux.Post("/envs/delete/:id", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.deleteEnv)))

	mux.Get("/oauth/redirect", http.HandlerFunc(app.oauthRedirect))
	mux.Get("/logout", http.HandlerFunc(app.logout))

	return standardMiddlware.Then(mux)
}
