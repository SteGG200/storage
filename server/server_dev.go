//go:build !production

package server

import (
	"net/http"

	"github.com/SteGG200/storage/server/middleware"
)

func setMiddleware(handler http.Handler) http.Handler {
	return middleware.Chain(handler,
		middleware.Log,
	)
}
