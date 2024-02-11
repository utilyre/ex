package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/utilyre/ex/router"
	"github.com/utilyre/xmate"
)

type loggerResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggerResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func NewLogger(logger *slog.Logger) router.Middleware {
	return func(next xmate.Handler) xmate.Handler {
		return xmate.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			w2 := &loggerResponseWriter{ResponseWriter: w}
			if err := next.ServeHTTP(w2, r); err != nil {
				logger.Warn("failed to run http handler",
					slog.String("remote", r.RemoteAddr),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("error", err.Error()),
				)

				return err
			}

			logger.Info("ran http handler",
				slog.String("remote", r.RemoteAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", w2.status),
			)

			return nil
		})
	}
}
