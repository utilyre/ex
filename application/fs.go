package application

import (
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(name string) (http.File, error) {
	f, err := nfs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if !s.IsDir() {
		return f, nil
	}

	if _, err = nfs.fs.Open(filepath.Join(name, "index.html")); err != nil {
		if err := f.Close(); err != nil {
			return nil, err
		}

		return nil, err
	}

	return f, nil
}
