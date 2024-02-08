package middlewares

import (
	"log/slog"
	"net/http"
)

func NewLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("running http handler",
				slog.String("remote", r.RemoteAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			next.ServeHTTP(w, r)
		})
	}
}
