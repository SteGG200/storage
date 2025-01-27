package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/config"
	"github.com/SteGG200/storage/server/mux"
	"github.com/SteGG200/storage/server/utils"
)

type Mux struct {
	*mux.Mux
}

func New(config *config.Config) http.Handler {
	router := &Mux{
		mux.New(config),
	}

	router.HandleFunc("/token/{path...}", router.sendToken)
	router.HandleFunc("/file", router.uploadData)

	return setMiddleware(router)
}

func (router *Mux) sendToken(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	filename := r.FormValue("name")

	err := db.RemoveUploadSessionByPath(router.Config.GetDatabase(), "/"+path, filename)

	if err != nil {
		logger.ErrorLogger.Print("Cannot remove old session")
	}

	data := map[string]any{
		"path":      "/" + path,
		"createdAt": time.Now().Unix(),
		"filename":  filename,
	}

	token, err := utils.GenerateToken(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["token"] = token

	err = db.SaveUploadSession(router.Config.GetDatabase(), data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = createFile(fmt.Sprintf("%s/%s", router.Config.GetStoragePath(), path), filename)

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

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"token": token,
	})
}

func (router *Mux) uploadData(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("name")
	file, header, _ := r.FormFile("file")
	token := r.FormValue("token")
	isLast := r.FormValue("isLast")

	session, err := utils.Verify(router.Config.GetDatabase(), token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if filename != session.Filename {
		http.Error(w, "Filename does not match", http.StatusBadRequest)
		return
	}

	path := strings.TrimPrefix(session.Path, "/")

	data := make([]byte, header.Size)

	readSize, err := file.Read(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if readSize < int(header.Size) {
		logger.WarningLogger.Print("Cannot read entire file.")
	}

	err = saveFile(data, fmt.Sprintf("%s/%s", router.Config.GetStoragePath(), path), filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isLast == "1" {
		err = db.RemoveUploadSessionByToken(router.Config.GetDatabase(), token)

		if err != nil {
			logger.ErrorLogger.Print(err)
		}
	}

	w.Write([]byte("OK"))
}
