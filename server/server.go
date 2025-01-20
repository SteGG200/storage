package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/config"
)

type Server struct {
	http   *http.Server
	config *config.Config
}

func New(db *string, configs ...config.ConfigFunc) (server *Server) {
	server = &Server{
		http: &http.Server{
			MaxHeaderBytes: 1 << 20, // 1MB
			WriteTimeout:   6 * time.Second,
			ReadTimeout:    6 * time.Second,
			IdleTimeout:    6 * time.Second,
		},
		config: config.New(db, configs...),
	}

	server.http.Handler = NewRouter(server.config)

	return
}

func (server *Server) Start(port string) error {
	server.http.Addr = fmt.Sprintf(":%s", port)
	logger.InfoLogger.Printf("Listening on port %s", port)
	return server.http.ListenAndServe()
}
