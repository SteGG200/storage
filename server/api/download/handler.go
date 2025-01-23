package download

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/SteGG200/storage/server/config"
	"github.com/SteGG200/storage/server/exception"
	"github.com/SteGG200/storage/server/mux"
)

type Mux struct {
	*mux.Mux
}

func New(config *config.Config) http.Handler {
	router := &Mux{
		mux.New(config),
	}

	router.HandleFunc("/{path...}", router.serveData)

	return router
}

func (router *Mux) serveData(w http.ResponseWriter, r *http.Request) {
	file, err := readFile(fmt.Sprintf("%s/%s", router.Config.GetStoragePath(), r.PathValue("path")))

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
}
