package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/utilyre/ex/router"
	"github.com/utilyre/xmate"
)

func NewRecoverer(logger *slog.Logger) router.Middleware {
	return func(next xmate.Handler) xmate.Handler {
		return xmate.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				msg := recover()
				if msg == nil {
					return
				}

				err = fmt.Errorf("panic: %v", msg)
			}()

			return next.ServeHTTP(w, r)
		})
	}
}
