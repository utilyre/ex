package router

import (
	"net/http"

	"github.com/utilyre/xmate"
)

type Middleware func(next xmate.Handler) xmate.Handler

type Router struct {
	mux          *http.ServeMux
	errorHandler xmate.ErrorHandler
	middlewares  []Middleware
}

func New(eh xmate.ErrorHandler) *Router {
	return &Router{
		mux:          http.NewServeMux(),
		errorHandler: eh,
		middlewares:  []Middleware{},
	}
}

func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) Handle(pattern string, handler xmate.Handler, middlewares ...Middleware) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	r.mux.Handle(pattern, r.errorHandler.Handle(handler))
}

func (r *Router) HandleFunc(pattern string, handler xmate.HandlerFunc, middlewares ...Middleware) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler).ServeHTTP
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler).ServeHTTP
	}

	r.mux.HandleFunc(pattern, r.errorHandler.HandleFunc(handler))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
