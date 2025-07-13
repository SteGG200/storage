package auth

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

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

	router.Handle("POST /setPassword/{path...}", router.setPassword())
	router.Handle("POST /login/{path...}", router.login())
	router.Handle("GET /check/{path...}", router.check())
	router.Handle("GET /checkNeedAuth/{path...}", router.checkNeedAuth())
	router.Handle("DELETE /removePassword/{path...}", router.removePassword())

	return setMiddleware(router)
}

func (router *Mux) setPassword() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "/" + r.PathValue("path")
		oldPassword := r.FormValue("oldPassword")
		password := r.FormValue("password")

		hashedPassword, err := utils.HashPassword(password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		doesExist, err := db.CheckIfPathHasAuth(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if doesExist {
			if oldPassword == "" {
				http.Error(w, "Require old password to change password", http.StatusBadRequest)
				return
			}
			hashedOldPassword, err := db.GetPasswordOfPath(router.Config.GetDatabase(), path)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			isMatch := utils.VerifyPassword(oldPassword, hashedOldPassword)

			if !isMatch {
				http.Error(w, "Old password is incorrect!", http.StatusForbidden)
				return
			}

			err = db.ChangePasswordOfPath(router.Config.GetDatabase(), path, hashedPassword)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err = db.CreateAuthForPath(router.Config.GetDatabase(), path, hashedPassword)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Write([]byte("OK"))
	})
}

func (router *Mux) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "/" + r.PathValue("path")
		password := r.FormValue("password")

		doesExist, err := db.CheckIfPathHasAuth(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the item needs to authenticate
		if !doesExist {
			http.Error(w, "Item doesn't need authentication", http.StatusContinue)
			return
		}

		// Verify the password is correct
		hashedPassword, err := db.GetPasswordOfPath(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isMatch := utils.VerifyPassword(password, hashedPassword)

		if !isMatch {
			http.Error(w, "Password is incorrect!", http.StatusForbidden)
			return
		}

		// Generate a new jwt token that is authenticated
		var data map[string]any

		_, token := utils.GetAuthorizationHeader(r)

		if token != "" {
			data, err = utils.ParseToken(token)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			data = make(map[string]any)
			data["path"] = make([]any, 0)
		}

		for _, currentPath := range data["path"].([]any) {
			currentPathToStr, ok := currentPath.(string)

			if !ok {
				http.Error(w, "There is an error on server!", http.StatusInternalServerError)
				return
			}

			if currentPathToStr == path {
				http.Error(w, "User already had access to this directory!", http.StatusContinue)
				return
			}
		}

		data["path"] = append(data["path"].([]any), path)

		token, err = utils.GenerateToken(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]any{
			"token": token,
		})
	})
}

func (router *Mux) check() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		typeToken, token := utils.GetAuthorizationHeader(r)

		if token != "" && typeToken != "Bearer" {
			http.Error(w, "Invalid Token", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		path := r.PathValue("path")

		allDirectories := append([]string{""}, strings.Split(path, "/")...)
		currentPath := "/"

		for _, directory := range allDirectories {
			currentPath = filepath.Join(currentPath, directory)
			doesNeedAuth, err := db.CheckIfPathHasAuth(router.Config.GetDatabase(), currentPath)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !doesNeedAuth {
				continue
			}

			isAuthenticated, err := utils.VerifyAuthorizationToken(token, currentPath)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !isAuthenticated {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]any{
					"unauthorizedPath": currentPath,
				})
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (router *Mux) checkNeedAuth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "/" + r.PathValue("path")

		doesNeedAuth, err := db.CheckIfPathHasAuth(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]any{
			"needAuth": doesNeedAuth,
		})
	})
}

func (router *Mux) removePassword() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "/" + r.PathValue("path")
		oldPassword := r.FormValue("oldPassword")

		doesNeedAuth, err := db.CheckIfPathHasAuth(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !doesNeedAuth {
			http.Error(w, "This directory isn't protected by password", http.StatusBadRequest)
			return
		}

		hashedOldPassword, err := db.GetPasswordOfPath(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isMatch := utils.VerifyPassword(oldPassword, hashedOldPassword)

		if !isMatch {
			http.Error(w, "Old password is incorrect!", http.StatusForbidden)
			return
		}

		err = db.RemovePasswordOfPath(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Remove password protection successfully"))
	})
}

func setMiddleware(handler http.Handler) http.Handler {
	return middleware.Chain(handler,
		validation,
	)
}
