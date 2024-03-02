package routes

import (
	"net/http"

	"github.com/utilyre/xmate/v2"
)

var ErrPageNotFound = xmate.Errorf(http.StatusNotFound, "Page Not Found")

type notFoundResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *notFoundResponseWriter) WriteHeader(status int) {
	w.status = status
	if status == http.StatusNotFound {
		return
	}

	w.ResponseWriter.WriteHeader(status)
}

func (w *notFoundResponseWriter) Write(p []byte) (int, error) {
	if w.status == http.StatusNotFound {
		return len(p), nil
	}

	return w.ResponseWriter.Write(p)
}

type Public struct {
	FileServer http.Handler
}

func (p Public) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	w2 := &notFoundResponseWriter{ResponseWriter: w}
	p.FileServer.ServeHTTP(w2, r)
	if w2.status == http.StatusNotFound {
		return ErrPageNotFound
	}

	return nil
}
