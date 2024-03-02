package application

import (
	"net/http"
	"path/filepath"

	"github.com/utilyre/ex/controller"
	"github.com/utilyre/ex/middleware"
)

func (app *Application) setupRoutes() {
	app.router.Use(middleware.NewLogger(app.logger))
	app.router.Use(middleware.NewRecoverer())

	home := controller.HomeController{
		HomeView: app.views.Lookup("home"),
	}
	app.router.HandleFunc("GET /{$}", home.Page)

	public := controller.PublicController{
		FileServer: http.FileServer(neuteredFileSystem{
			fs: http.Dir(filepath.Join(app.cfg.AppRoot, "public")),
		}),
	}
	app.router.Handle("/", public)
}
