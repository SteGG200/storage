package auth

import (
	"net/http"
	"strings"
)

func validation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/check") {
			next.ServeHTTP(w, r)
			return
		}

		password := r.FormValue("password")

		if password == "" {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
