package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/config"
)

type Server struct {
	http   *http.Server
	config *config.Config
}

/*
New creates a new instance of Server with the provided database and optional configuration functions.

Params:

	  db *db.DB // The database instance.
		configs...config.ConfigFunc // Optional configuration functions.

Returns:

	*Server // A new instance of Server with the provided database and optional configuration functions.
*/
func New(db *db.DB, configs ...config.ConfigFunc) (server *Server) {
	server = &Server{
		http: &http.Server{
			MaxHeaderBytes: 1 << 20, // 1MB
			WriteTimeout:   6 * time.Second,
			ReadTimeout:    6 * time.Second,
			IdleTimeout:    6 * time.Second,
		},
		config: config.New(db, configs...),
	}

	server.http.Handler = setMiddleware(NewRouter(server.config))

	return
}

/*
Start starts the server and listens on the specified port.

Params:

	port string // The port to listen on.
*/
func (server *Server) Start(port string) error {
	server.http.Addr = fmt.Sprintf(":%s", port)
	logger.InfoLogger.Printf("Listening on port %s", port)
	return server.http.ListenAndServe()
}
