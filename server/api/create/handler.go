package create

import (
	"net/http"

	"github.com/SteGG200/storage/server/config"
	"github.com/SteGG200/storage/server/mux"
)

type Mux struct {
	*mux.Mux
}

func New(config *config.Config) http.Handler {
	router := &Mux{
		mux.New(config),
	}

	router.HandleFunc("/{path...}", router.createFolder)

	return router
}

func (router *Mux) createFolder(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	foldername := r.FormValue("name")

	err := createFolder(router.Config.GetStoragePath()+"/"+path, foldername)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Folder created successfully"))
}
