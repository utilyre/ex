package application

import (
	"database/sql"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/utilyre/ex/config"
	"github.com/utilyre/ex/middlewares"
	"github.com/utilyre/ex/routes"
	"github.com/utilyre/xmate"
)

type Application struct {
	cfg      config.Config
	logger   *slog.Logger
	views    *template.Template
	router   chi.Router
	handler  xmate.ErrorHandler
	validate *validator.Validate
	db       *bun.DB
}

func New(cfg config.Config, logger *slog.Logger) *Application {
	views, err := template.ParseGlob(filepath.Join(cfg.AppRoot, "views", "*.html"))
	if err != nil {
		logger.Error("failed to parse views", "error", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	handler := newHandler(logger, views.Lookup("error"))
	validate := validator.New()

	sqldb, err := sql.Open(sqliteshim.ShimName, cfg.DSN)
	if err != nil {
		logger.Error("failed to open connection to database", "dsn", cfg.DSN, "error", err)
		os.Exit(1)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())

	return &Application{
		cfg:      cfg,
		logger:   logger,
		views:    views,
		router:   router,
		handler:  handler,
		validate: validate,
		db:       db,
	}
}

func (app *Application) Setup() *Application {
	app.router.Use(
		middlewares.NewRecoverer(app.logger),
		middlewares.NewLogger(app.logger),
	)

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

func newHandler(logger *slog.Logger, errorView *template.Template) xmate.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.Context().Value(xmate.ErrorKey{}).(error)

		httpErr := new(xmate.HTTPError)
		if !errors.As(err, &httpErr) {
			httpErr.Code = http.StatusInternalServerError
			httpErr.Message = "Internal Server Error"

			logger.Warn("failed to run http handler",
				slog.String("remote", r.RemoteAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("error", err.Error()),
			)
		}

		_ = xmate.WriteHTML(w, errorView, httpErr.Code, httpErr)
	}
}