//go:build production

package middleware

import (
	"net/http"
)

// Logger is a no-op in production.
func Logger(next http.Handler) http.Handler {
	return next
}
