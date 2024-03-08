package router

import (
	"net/http"

	"github.com/utilyre/xmate/v2"
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
	handler = wrapHandler(handler, middlewares)
	handler = wrapHandler(handler, r.middlewares)

	r.serveMux.Handle(pattern, r.handler.Handle(handler))
}

func (r *Router) HandleFunc(pattern string, handler xmate.HandlerFunc, middlewares ...Middleware) {
	handler = wrapHandler(handler, middlewares).ServeHTTP
	handler = wrapHandler(handler, r.middlewares).ServeHTTP

	r.serveMux.HandleFunc(pattern, r.handler.HandleFunc(handler))
}

func wrapHandler(handler xmate.Handler, middlewares []Middleware) xmate.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
