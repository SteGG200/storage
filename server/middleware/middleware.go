package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

// Chain combines multiple middleware functions into a single middleware handler.
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}
