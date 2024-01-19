package routes

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/utilyre/xmate"
)

type Home struct {
	Handler  xmate.ErrorHandler
	HomeView *template.Template
}

func (h Home) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.Handler.HandleFunc(h.page))

	return r
}

func (h Home) page(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, h.HomeView, http.StatusOK, nil)
}
