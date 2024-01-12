package application

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/utilyre/xmate"
)

type Application struct {
	logger  *slog.Logger
	router  chi.Router
	handler xmate.ErrorHandler
	views   *template.Template
}

func New() *Application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	views, err := template.ParseGlob(filepath.Join("views", "*.html"))
	if err != nil {
		logger.Error("failed to parse views", "error", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := r.Context().Value(xmate.ErrorKey{}).(error)

		if httpErr := new(xmate.HTTPError); errors.As(err, &httpErr) {
			_ = xmate.WriteText(w, httpErr.Code, httpErr.Message)
			return
		}

		logger.Warn("failed to run http handler",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("error", err.Error()),
		)

		_ = xmate.WriteText(w, http.StatusInternalServerError, "Internal Server Error")
	}

	return &Application{
		logger:  logger,
		router:  router,
		handler: handler,
		views:   views,
	}
}

func (app *Application) Init() {
	app.router.Get("/",
		app.handler.HandleFunc(func(w http.ResponseWriter, r *http.Request) error {
			return xmate.WriteHTML(w, app.views.Lookup("home"), http.StatusOK, nil)
		}),
	)

	// TODO: add routes

	app.router.Handle("/*", http.FileServer(neuteredFileSystem{
		fs: http.Dir("public"),
	}))
}

func (app *Application) Start() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.router,
	}

	app.logger.Info("starting to listen and serve", "address", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error("failed to listen and serve", "error", err)
		os.Exit(1)
	}
}
