package remove

import (
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/server/config"
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

	router.Handle("/{path...}", middleware.SetAuthorization(router.removeItem(), "path", router.Config.GetDatabase()))

	return router
}

func (router *Mux) removeItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		root := router.Config.GetStoragePath()
		path := "/" + r.PathValue("path")

		if path == "" {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		err := removeItem(filepath.Join(root, path))

		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				http.Error(w, "Not found directory", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db.RemovePasswordOfPath(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Deleted successfully"))
	})
}
