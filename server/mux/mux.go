package mux

import (
	"net/http"

	"github.com/SteGG200/storage/config"
)

type Mux struct {
	router *http.ServeMux
	Config config.Config
}

func New(config config.Config) (mux *Mux) {
	mux = &Mux{
		router: http.NewServeMux(),
		Config: config,
	}

	return
}

func (mux *Mux) Handle(pattern string, handler http.Handler) {
	mux.router.Handle(pattern, handler)
}

func (mux *Mux) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	mux.router.HandleFunc(pattern, handlerFunc)
}
