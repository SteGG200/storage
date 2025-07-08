package middleware

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/server/utils"
)

/*
SetAuthorization middleware sets the authorization for specific paths.

Params:

	next http.Handler // The next handler in the middleware chain.
	pathValue string   // The path value to extract the authorization token from.
	database *db.DB   // The database instance.
*/
func SetAuthorization(next http.Handler, pathValue string, database *db.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		typeToken, token := utils.GetAuthorizationHeader(r)

		if token != "" && typeToken != "Bearer" {
			http.Error(w, "Invalid Token", http.StatusBadRequest)
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
