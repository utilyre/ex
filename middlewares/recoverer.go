package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/utilyre/ex/router"
	"github.com/utilyre/xmate"
)

func NewRecoverer(logger *slog.Logger) router.Middleware {
	return func(next xmate.Handler) xmate.Handler {
		return xmate.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			defer func() {
				msg := recover()
				if msg == nil {
					return
				}

				logger.Warn("failed to run http handler (panicked)",
					slog.String("remote", r.RemoteAddr),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Any("message", msg),
				)
			}()

			return next.ServeHTTP(w, r)
		})
	}
}
