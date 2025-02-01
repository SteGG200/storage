package middleware

import (
	"net/http"

	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/server/utils"
)

func SetAuthorization(next http.Handler, pathValue string, database *db.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "/" + r.PathValue(pathValue)

		doesAuthNeed, err := db.CheckIfPathHasAuth(database, path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !doesAuthNeed {
			next.ServeHTTP(w, r)
			return
		}

		typeToken, token := utils.GetAuthorizationHeader(r)

		if typeToken != "Bearer" || token == "" {
			http.Error(w, "Invalid Authorization", http.StatusUnauthorized)
			return
		}

		isAuth, err := utils.VerifyAuthorizationToken(token, path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isAuth {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
