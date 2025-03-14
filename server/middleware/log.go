package middleware

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/SteGG200/storage/logger"
)

const MAX_ID = 10000

// SetLog middleware logs the incoming requests and their corresponding response times.
func SetLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := rand.Intn(MAX_ID)
		logger.InfoLogger.Printf("%s: %s (id: %d)", r.Method, r.URL.Path, id)
		startTime := time.Now()
		next.ServeHTTP(w, r)
		elapsedTime := time.Since(startTime)
		logger.InfoLogger.Printf("Response %d in %s", id, elapsedTime)
	})
}
