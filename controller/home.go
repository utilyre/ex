package controller

import (
	"html/template"
	"net/http"

	"github.com/utilyre/xmate/v2"
)

type HomeController struct {
	HomeView *template.Template
}

func (hc HomeController) Page(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, hc.HomeView, http.StatusOK, nil)
}
