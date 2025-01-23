package middleware

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/SteGG200/storage/logger"
)

const MAX_ID = 10000

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := rand.Intn(MAX_ID)
		logger.InfoLogger.Printf("%s: %s (id: %d)", r.Method, r.URL.Path, id)
		logger.InfoLogger.Print(r.Header)
		startTime := time.Now()
		next.ServeHTTP(w, r)
		elapsedTime := time.Since(startTime)
		logger.InfoLogger.Print(w.Header())
		logger.InfoLogger.Printf("Response %d in %s", id, elapsedTime)
	})
}
