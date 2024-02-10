package routes

import (
	"html/template"
	"net/http"

	"github.com/utilyre/xmate"
)

type Home struct {
	HomeView *template.Template
}

func (h Home) Page(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, h.HomeView, http.StatusOK, nil)
}
