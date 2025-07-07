package download

import (
	"errors"
	"io/fs"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/SteGG200/storage/server/config"
	"github.com/SteGG200/storage/server/exception"
	"github.com/SteGG200/storage/server/middleware"
	"github.com/SteGG200/storage/server/mux"
)

type Mux struct {
	*mux.Mux
}

func New(config *config.Config) http.Handler {
	router := &Mux{
		mux.New(config),
	}

	router.Handle("/{path...}", middleware.SetAuthorization(router.serveData(), "path", router.Config.GetDatabase()))

	return router
}

func (router *Mux) serveData() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		root := router.Config.GetStoragePath()
		path, _ := url.PathUnescape(r.PathValue("path"))
		file, err := readFile(filepath.Join(root, path))

		if err != nil {
			if errors.Is(err, fs.ErrNotExist) || errors.Is(err, exception.ErrNotAFile) {
				http.Error(w, "Not found file", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		stat, err := file.Stat()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
	})
}
