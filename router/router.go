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

func (r *Router) Group(fn func(r *Router)) {
	fn(&Router{
		mux:          r.mux,
		errorHandler: r.errorHandler,
		middlewares:  append([]Middleware(nil), r.middlewares...),
	})
}

func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) HandleUnsafe(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

func (r *Router) Handle(pattern string, handler xmate.Handler) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	r.mux.Handle(pattern, r.errorHandler.Handle(handler))
}

func (r *Router) HandleFunc(pattern string, handler xmate.HandlerFunc) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler).ServeHTTP
	}

	r.mux.HandleFunc(pattern, r.errorHandler.HandleFunc(handler))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
