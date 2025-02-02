package auth

import (
	"encoding/json"
	"net/http"
	"slices"

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

	router.Handle("/register/{path...}", router.register())
	router.Handle("/login/{path...}", router.login())

	return setMiddleware(router)
}

func (router *Mux) register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "/" + r.PathValue("path")
		password := r.FormValue("password")

		doesExist, err := db.CheckIfPathHasAuth(router.Config.GetDatabase(), path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if doesExist {
			http.Error(w, "Item already has authentication", http.StatusConflict)
			return
		}

		hashedPassword, err := utils.HashPassword(password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db.CreateAuthForPath(router.Config.GetDatabase(), path, hashedPassword)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

		if !doesExist {
			http.Error(w, "Item doesn't need authentication", http.StatusContinue)
			return
		}

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
			data["path"] = make([]string, 0)
		}

		if slices.Contains(data["path"].([]string), path) {
			http.Error(w, "Item is already authenticated", http.StatusContinue)
			return
		}

		data["path"] = append(data["path"].([]string), path)

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

func setMiddleware(handler http.Handler) http.Handler {
	return middleware.Chain(handler)
}
