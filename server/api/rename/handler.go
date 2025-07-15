package rename

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

	router.Handle("/{path...}", middleware.SetAuthorization(router.renameItem(), "path", router.Config.GetDatabase()))

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

		// Update folder's name in database
		err = db.RenameAuthPath(router.Config.GetDatabase(), path, filepath.Join(filepath.Dir(path), newName))

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

func setMiddleware(handler http.Handler) http.Handler {
	return middleware.Chain(handler,
		validation,
	)
}
