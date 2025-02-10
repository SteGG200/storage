package rename

import "net/http"

func validation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newName := r.FormValue("newName")

		if newName == "" {
			http.Error(w, "New name is required", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
