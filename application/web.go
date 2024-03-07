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

	setupHomeController(app)
	setupPublicController(app)
}

func setupHomeController(app *Application) {
	hc := controller.HomeController{
		HomeView: app.views.Lookup("home"),
	}

	app.router.HandleFunc("GET /{$}", hc.Page)
}

func setupPublicController(app *Application) {
	pc := controller.PublicController{
		FileServer: http.FileServer(neuteredFileSystem{
			fs: http.Dir(filepath.Join(app.cfg.AppRoot, "public")),
		}),
	}

	app.router.Handle("/", pc)

}
