package application

import (
	"net/http"
	"path/filepath"

	"github.com/utilyre/ex/controller"
	"github.com/utilyre/ex/middleware"
)

func (app *Application) setupWeb() {
	app.router.Use(middleware.NewLogger(app.logger))
	app.router.Use(middleware.NewRecoverer())

	homeController := newHomeController(app)
	app.router.HandleFunc("GET /{$}", homeController.Page)

	publicController := newPublicController(app)
	app.router.Handle("/", publicController)
}

func newHomeController(app *Application) controller.HomeController {
	return controller.HomeController{
		HomeView: app.views.Lookup("home"),
	}
}

func newPublicController(app *Application) controller.PublicController {
	return controller.PublicController{
		FileServer: http.FileServer(neuteredFileSystem{
			fs: http.Dir(filepath.Join(app.cfg.AppRoot, "public")),
		}),
	}
}
