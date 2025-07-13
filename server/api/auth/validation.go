package auth

import (
	"net/http"
	"strings"
)

func validation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/removeAuth") {
			oldPassword := r.FormValue("oldPassword")

			if oldPassword == "" {
				http.Error(w, "Old password is required", http.StatusBadRequest)
				return
			}
		}

		if strings.HasPrefix(r.URL.Path, "/setPassword") || strings.HasPrefix(r.URL.Path, "/login") {
			password := r.FormValue("password")

			if password == "" {
				http.Error(w, "Password is required", http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
