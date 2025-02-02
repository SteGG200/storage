package middleware

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/server/utils"
)

func SetAuthorization(next http.Handler, pathValue string, database *db.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		typeToken, token := utils.GetAuthorizationHeader(r)

		if typeToken != "Bearer" && token != "" {
			http.Error(w, "Invalid Authorization", http.StatusUnauthorized)
			return
		}

		path := r.PathValue(pathValue)

		allDirectories := append([]string{""}, strings.Split(path, "/")...)
		currentPath := "/"

		for _, directory := range allDirectories {
			currentPath = filepath.Join(currentPath, directory)
			doesNeedAuth, err := db.CheckIfPathHasAuth(database, currentPath)

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
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
