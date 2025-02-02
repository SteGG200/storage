package auth

import "net/http"

func validation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := r.FormValue("password")

		if password == "" {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
