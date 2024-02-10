package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/utilyre/ex/router"
	"github.com/utilyre/xmate"
)

func NewLogger(logger *slog.Logger) router.Middleware {
	return func(next xmate.Handler) xmate.Handler {
		return xmate.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			logger.Info("running http handler",
				slog.String("remote", r.RemoteAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			return next.ServeHTTP(w, r)
		})
	}
}
