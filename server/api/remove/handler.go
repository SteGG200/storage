package remove

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/SteGG200/storage/db"
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

		// Remove password of removed folder from database
		err = db.RemovePasswordOfPath(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Remove the removed path from data in jwt token
		w.Header().Set("Content-Type", "application/json")

		_, token := utils.GetAuthorizationHeader(r)

		if token != "" {
			data, err := utils.ParseToken(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			newData := make(map[string]any)
			newData["path"] = make([]any, 0)
			for _, currentPath := range data["path"].([]any) {
				currentPathToStr, ok := currentPath.(string)

				if !ok {
					http.Error(w, "Unexpected type of path in data of jwt token", http.StatusInternalServerError)
					return
				}

				if currentPathToStr == path {
					continue
				}

				newData["path"] = append(newData["path"].([]any), currentPath)
			}

			newToken, err := utils.GenerateToken(newData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]any{
				"token": newToken,
			})

			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"token": "",
		})
	})
}
