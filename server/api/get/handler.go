package get

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
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
		path := r.PathValue("path")
		entries, err := listEntries(filepath.Join(root, path))

		if err != nil {
			if errors.Is(err, fs.ErrNotExist) || errors.Is(err, exception.ErrNotADirectory) {
				http.Error(w, "Not found directory", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(entries)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})
}
