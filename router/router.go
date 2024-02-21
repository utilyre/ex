package router

import (
	"net/http"

	"github.com/utilyre/xmate"
)

type serveMux = http.ServeMux
type Middleware func(next xmate.Handler) xmate.Handler

type Router struct {
	*serveMux
	handler     xmate.ErrorHandler
	middlewares []Middleware
}

func New(handler xmate.ErrorHandler) *Router {
	return &Router{
		serveMux:    http.NewServeMux(),
		handler:     handler,
		middlewares: []Middleware{},
	}
}

func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) Handle(pattern string, handler xmate.Handler, middlewares ...Middleware) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	r.serveMux.Handle(pattern, r.handler.Handle(handler))
}

func (r *Router) HandleFunc(pattern string, handler xmate.HandlerFunc, middlewares ...Middleware) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler).ServeHTTP
	}
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler).ServeHTTP
	}

	r.serveMux.HandleFunc(pattern, r.handler.HandleFunc(handler))
}
