package upload

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/config"
	"github.com/SteGG200/storage/server/middleware"
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

	router.Handle("/token/{path...}", middleware.SetAuthorization(router.sendToken(), "path", router.Config.GetDatabase()))
	router.Handle("/file", router.uploadData())

	return setMiddleware(router)
}

func (router *Mux) sendToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		root := router.Config.GetStoragePath()
		path := "/" + r.PathValue("path")
		filename := r.FormValue("name")

		err := db.RemoveUploadSessionByPath(router.Config.GetDatabase(), path, filename)

		if err != nil {
			logger.ErrorLogger.Print("Cannot remove old session")
		}

		data := map[string]any{
			"path":      path,
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

		err = createFile(filepath.Join(root, path), filename)

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
	})
}

func (router *Mux) uploadData() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		root := router.Config.GetStoragePath()
		filename := r.FormValue("name")
		file, header, _ := r.FormFile("file")
		_, token := utils.GetAuthorizationHeader(r)
		isLast := r.FormValue("isLast")

		session, err := utils.VerifyUploadSession(router.Config.GetDatabase(), token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if filename != session.Filename {
			http.Error(w, "Filename does not match", http.StatusBadRequest)
			return
		}

		data := make([]byte, header.Size)

		readSize, err := file.Read(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if readSize < int(header.Size) {
			logger.WarningLogger.Print("Cannot read entire file.")
		}

		err = saveFile(data, filepath.Join(root, session.Path), filename)

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
	})
}

func setMiddleware(handler http.Handler) http.Handler {
	return middleware.Chain(handler,
		validation,
		preflightHandler,
	)
}
