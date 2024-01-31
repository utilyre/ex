package application

import (
	"database/sql"
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/utilyre/ex/config"
	"github.com/utilyre/ex/routes"
	"github.com/utilyre/xmate"
)

type Application struct {
	cfg     config.Config
	logger  *slog.Logger
	views   *template.Template
	router  chi.Router
	handler xmate.ErrorHandler
	db      *bun.DB
}

func New(cfg config.Config) *Application {
	logger := newLogger(cfg)

	views, err := template.ParseGlob(filepath.Join(cfg.AppRoot, "views", "*.html"))
	if err != nil {
		logger.Error("failed to parse views", "error", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	handler := newHandler(logger, views.Lookup("error"))

	logger.Info("opening connection to database", "dsn", cfg.DSN)
	sqldb, err := sql.Open(sqliteshim.ShimName, cfg.DSN)
	if err != nil {
		logger.Error("failed to open connection to database", "error", err)
		os.Exit(1)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())

	return &Application{
		cfg:     cfg,
		logger:  logger,
		views:   views,
		router:  router,
		handler: handler,
		db:      db,
	}
}

func (app *Application) Setup() *Application {
	app.router.Mount("/assets", http.StripPrefix(
		"/assets",
		http.FileServer(neuteredFileSystem{
			fs: http.Dir(filepath.Join(app.cfg.AppRoot, "assets")),
		}),
	))

	app.router.Mount("/", routes.Home{
		Handler:  app.handler,
		HomeView: app.views.Lookup("home"),
	}.Router())

	return app
}

func (app *Application) Start() {
	app.logger.Info("starting server application", "config", app.cfg)

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
		f, err := os.OpenFile(filepath.Join(cfg.AppRoot, "app.log"),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}

		handler = slog.NewJSONHandler(io.MultiWriter(os.Stdout, f), opts)
	}

	return slog.New(handler)
}

func newHandler(logger *slog.Logger, errorView *template.Template) xmate.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.Context().Value(xmate.ErrorKey{}).(error)

		httpErr := new(xmate.HTTPError)
		if !errors.As(err, &httpErr) {
			httpErr.Code = http.StatusInternalServerError
			httpErr.Message = "Internal Server Error"

			logger.Warn("failed to run http handler",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("error", err.Error()),
			)
		}

		_ = xmate.WriteHTML(w, errorView, httpErr.Code, httpErr)
	}
}
