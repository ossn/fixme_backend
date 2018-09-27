package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/envy"
	tokenauth "github.com/gobuffalo/mw-tokenauth"
	"github.com/gobuffalo/x/sessions"
	"github.com/ossn/fixme_backend/models"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {

	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			Prefix:      "/api",
			SessionName: "_fixme_backend_session",
		})
		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		app.GET("/projects", ProjectsResource{}.List)
		app.GET("/repositories", RepositoriesResource{}.List)
		app.GET("/issues", IssuesResource{}.ListOpen)
		app.GET("/issues-count", IssuesResource{}.Count)
		app.POST("/login", AdminsResource{}.Login)

		admin := app.Group("/admin")
		admin.Use(tokenauth.New(tokenauth.Options{}))

		admin.Resource("/projects", ProjectsResource{})
		admin.Resource("/repositories", RepositoriesResource{})
		admin.Resource("/issues", IssuesResource{})
		admin.Resource("/users", AdminsResource{})
	}
	return app
}
