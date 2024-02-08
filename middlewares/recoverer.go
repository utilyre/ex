package middlewares

import (
	"log/slog"
	"net/http"
)

func NewRecoverer(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				msg := recover()
				if msg == nil {
					return
				}

				logger.Warn("failed to run http handler (panicked)",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Any("message", msg),
				)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
