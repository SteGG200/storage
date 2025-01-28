package upload

import (
	"net/http"
	"strings"

	"github.com/SteGG200/storage/server/utils"
)

func validation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("name")
		if filename == "" {
			http.Error(w, "Filename is required", http.StatusBadRequest)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/token") {
			next.ServeHTTP(w, r)
			return
		}

		file, header, err := r.FormFile("file")

		if err != nil || header == nil || file == nil {
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}

		typeToken, token := utils.GetAuthorizationHeader(r)

		if typeToken != "Bearer" || token == "" {
			http.Error(w, "Invalid Authorization", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
