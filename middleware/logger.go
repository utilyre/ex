package middleware

import (
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/utilyre/ex/application/router"
	"github.com/utilyre/xmate/v2"
)

func NewLogger(logger *slog.Logger) router.Middleware {
	return func(next xmate.Handler) xmate.Handler {
		return xmate.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				return err
			}

			w2 := &loggerResponseWriter{ResponseWriter: w}
			if err := next.ServeHTTP(w2, r); err != nil {
				if httpErr := (xmate.HTTPError{}); errors.As(err, &httpErr) {
					w2.status = httpErr.Code
				} else {
					w2.status = http.StatusInternalServerError
				}

				logger.Warn("failed to run http handler",
					slog.String("ip", ip),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", w2.status),
					slog.String("error", err.Error()),
				)

				return err
			}

			logger.Info("ran http handler",
				slog.String("ip", ip),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", w2.status),
			)

			return nil
		})
	}
}

type loggerResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggerResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
