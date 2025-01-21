package upload

import (
	"net/http"

	"github.com/SteGG200/storage/server/middleware"
)

func validation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("name")
		if filename == "" {
			http.Error(w, "Filename is required", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")

		if err != nil || header == nil || file == nil {
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func setMiddleware(handler http.Handler) http.Handler {
	return middleware.Chain(handler,
		validation,
	)
}
