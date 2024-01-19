package application

import (
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/utilyre/ex/config"
	"github.com/utilyre/ex/routes"
	"github.com/utilyre/xmate"
)

type Application struct {
	cfg     config.Config
	logger  *slog.Logger
	router  chi.Router
	handler xmate.ErrorHandler
	views   *template.Template
}

func New(cfg config.Config) *Application {
	logger := newLogger(cfg)

	router := chi.NewRouter()
	handler := newHandler(logger)

	views, err := template.ParseGlob(filepath.Join(cfg.Root, "views", "*.html"))
	if err != nil {
		logger.Error("failed to parse views", "error", err)
		os.Exit(1)
	}

	return &Application{
		cfg:     cfg,
		logger:  logger,
		router:  router,
		handler: handler,
		views:   views,
	}
}

func (app *Application) Init() {
	app.router.Mount("/public", http.StripPrefix(
		"/public",
		http.FileServer(neuteredFileSystem{
			fs: http.Dir(filepath.Join(app.cfg.Root, "public")),
		}),
	))

	app.router.Mount("/", routes.Home{
		Handler:  app.handler,
		HomeView: app.views.Lookup("home"),
	}.Router())
}

func (app *Application) Start() {
	srv := &http.Server{
		Addr:    app.cfg.ServerAddr,
		Handler: app.router,
	}

	app.logger.Info("starting to listen and serve", "address", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error("failed to listen and serve", "error", err)
		os.Exit(1)
	}
}

func newLogger(cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}

	var handler slog.Handler
	switch cfg.Mode {
	case config.ModeDev:
		handler = slog.NewTextHandler(os.Stdout, opts)
	case config.ModeProd:
		f, err := os.OpenFile(filepath.Join(cfg.Root, "app.log"),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}

		handler = slog.NewJSONHandler(io.MultiWriter(os.Stdout, f), opts)
	}

	return slog.New(handler)
}

func newHandler(logger *slog.Logger) xmate.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) {
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

		_ = xmate.WriteText(w,
			http.StatusInternalServerError, "Internal Server Error")
	}
}
