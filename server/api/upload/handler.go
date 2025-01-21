package upload

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/SteGG200/storage/logger"
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

	router.HandleFunc("/{path...}", router.serveData)

	return setMiddleware(router)
}

func (router *Mux) serveData(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("name")
	file, header, _ := r.FormFile("file")

	data := make([]byte, header.Size)

	readSize, err := file.Read(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if readSize < int(header.Size) {
		logger.WarningLogger.Print("Cannot read entire file.")
	}

	err = saveFile(data, fmt.Sprintf("%s/%s", router.Config.GetStoragePath(), r.PathValue("path")), filename)

	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			http.Error(w, "File already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "Directory does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}
