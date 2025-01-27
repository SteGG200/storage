package create

import "net/http"

func validation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		foldername := r.FormValue("name")

		if foldername == "" {
			http.Error(w, "Folder name is required", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
