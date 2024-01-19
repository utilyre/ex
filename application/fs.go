package application

import (
	"net/http"
	"os"
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
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}
