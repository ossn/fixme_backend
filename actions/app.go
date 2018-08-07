package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/envy"

	"github.com/gobuffalo/x/sessions"
	"github.com/ossn/fixme_backend/models"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

var Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

// QOR admin
// var Admin = admin.New(&qor.Config{DB: &gorm.DB{}})

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.

func App() *buffalo.App {
	// Admin.AddResource(&models.Project{})
	// Admin.AddResource(&models.Issue{})
	// Admin.AddResource(&models.Repository{})
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
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

		app.Resource("/projects", ProjectsResource{})
		app.Resource("/repositories", RepositoriesResource{})
		app.Resource("/issues", IssuesResource{})
		app.GET("/issues-count", IssuesResource{}.Count)
		// app.ANY("/admin/{path:.+}", buffalo.WrapHandler(http.StripPrefix("/admin", Other())))
	}
	return app
}

// func Other() http.Handler {
// 	f := func(res http.ResponseWriter, req *http.Request) {
// 		fmt.Fprintln(res, req.URL.String())
// 		fmt.Fprintln(res, req.Method)
// 	}
// 	mux := http.NewServeMux()
// 	// mux := mux.NewRouter()
// 	Admin.MountTo("/admin", mux)
// 	return mux
// }
