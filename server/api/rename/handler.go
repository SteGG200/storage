package rename

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

	router.Handle("/{path...}", router.renameItem())

	return setMiddleware(router)
}

func (router *Mux) renameItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		root := router.Config.GetStoragePath()
		path := "/" + r.PathValue("path")
		newName := r.FormValue("newName")

		err := renameItem(filepath.Join(root, path), newName)

		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				http.Error(w, "Not found directory", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db.RenameAuthPath(router.Config.GetDatabase(), path, filepath.Join(filepath.Dir(path), newName))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Renamed successfully"))
	})
}

func setMiddleware(handler http.Handler) http.Handler {
	return middleware.Chain(handler,
		validation,
	)
}
