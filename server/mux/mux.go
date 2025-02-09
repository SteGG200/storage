package mux

import (
	"net/http"

	"github.com/SteGG200/storage/server/config"
)

type Mux struct {
	router *http.ServeMux
	Config *config.Config
}

/*
New creates a new instance of Mux with the provided configuration.

Params:

	config *config.Config // The server configuration.

Returns:

	*Mux // A new instance of Mux with the provided configuration.
*/
func New(config *config.Config) (mux *Mux) {
	mux = &Mux{
		router: http.NewServeMux(),
		Config: config,
	}

	return
}

/*
Handle attaches a handler to the given URL pattern.

Params:

	pattern string // The URL pattern to match.
	handler http.Handler // The handler function to execute for the matching URL pattern.

Note: This function is a convenience method, equivalent to calling HandleFunc with a single function argument.
*/
func (mux *Mux) Handle(pattern string, handler http.Handler) {
	mux.router.Handle(pattern, handler)
}

/*
HandleFunc attaches a handler function to the given URL pattern.

Params:

	pattern string // The URL pattern to match.
	handlerFunc http.HandlerFunc // The handler function to execute for the matching URL pattern.
*/
func (mux *Mux) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	mux.router.HandleFunc(pattern, handlerFunc)
}

/*
ServeHTTP dispatches the request to the handler whose URL path matches the request's URL path.

Params:

	w http.ResponseWriter // The HTTP response writer.
	r *http.Request // The HTTP request.

Note: This function is required for http.Handler interface.
*/
func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.router.ServeHTTP(w, r)
}
