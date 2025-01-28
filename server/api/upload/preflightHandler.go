package upload

import (
	"net/http"

	"github.com/SteGG200/storage/logger"
)

func preflightHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			logger.InfoLogger.Print("response OPTIONS request")
			return
		}

		next.ServeHTTP(w, r)
	})
}
