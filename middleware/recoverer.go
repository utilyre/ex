package middleware

import (
	"fmt"
	"net/http"

	"github.com/utilyre/ex/router"
	"github.com/utilyre/xmate/v2"
)

func NewRecoverer() router.Middleware {
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
