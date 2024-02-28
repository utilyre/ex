package application

import (
	"database/sql"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/utilyre/ex/config"
	"github.com/utilyre/ex/router"
	"github.com/utilyre/xmate"
)

type Application struct {
	cfg      config.Config
	logger   *slog.Logger
	views    *template.Template
	router   *router.Router
	validate *validator.Validate
	db       *bun.DB
}

func New(cfg config.Config) *Application {
	logger := newLogger(cfg.Mode, cfg.LogLevel)

	views, err := template.ParseGlob(filepath.Join(cfg.AppRoot, "views", "*.html"))
	if err != nil {
		logger.Error("failed to parse views", "error", err)
		os.Exit(1)
	}

	router := router.New(newErrorHandler(views.Lookup("error")))
	validate := validator.New()

	sqldb, err := sql.Open(sqliteshim.ShimName, cfg.DSN)
	if err != nil {
		logger.Error("failed to open database connection", "dsn", cfg.DSN, "error", err)
		os.Exit(1)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())

	return &Application{
		cfg:      cfg,
		logger:   logger,
		views:    views,
		router:   router,
		validate: validate,
		db:       db,
	}
}

func (app *Application) Setup() *Application {
	app.setupRoutes()
	return app
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

func newLogger(mode config.Mode, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	switch mode {
	case config.ModeDev:
		handler = slog.NewTextHandler(os.Stdout, opts)
	case config.ModeProd:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func newErrorHandler(errorView *template.Template) xmate.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.Context().Value(xmate.ErrorKey{}).(error)

		httpErr := new(xmate.HTTPError)
		if !errors.As(err, &httpErr) {
			httpErr.Code = http.StatusInternalServerError
			httpErr.Message = "Internal Server Error"
		}

		_ = xmate.WriteHTML(w, errorView, httpErr.Code, httpErr)
	}
}
