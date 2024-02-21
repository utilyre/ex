package application

import (
	"net/http"
	"path/filepath"

	"github.com/utilyre/ex/middlewares"
	"github.com/utilyre/ex/routes"
)

func (app *Application) setupRoutes() {
	app.router.Use(middlewares.NewLogger(app.logger))
	app.router.Use(middlewares.NewRecoverer())

	home := routes.Home{
		HomeView: app.views.Lookup("home"),
	}
	app.router.HandleFunc("GET /{$}", home.Page)

	public := routes.Public{
		FileServer: http.FileServer(neuteredFileSystem{
			fs: http.Dir(filepath.Join(app.cfg.AppRoot, "public")),
		}),
	}
	app.router.Handle("/", public)
}
