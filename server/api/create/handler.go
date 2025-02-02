package create

import (
	"net/http"
	"path/filepath"

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

	router.Handle("/{path...}", middleware.SetAuthorization(router.createFolder(), "path", router.Config.GetDatabase()))

	return setMiddleware(router)
}

func (router *Mux) createFolder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		root := router.Config.GetStoragePath()
		path := r.PathValue("path")
		foldername := r.FormValue("name")

		err := createFolder(filepath.Join(root, path), foldername)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Folder created successfully"))
	})
}

func setMiddleware(next http.Handler) http.Handler {
	return middleware.Chain(next,
		validation,
	)
}
