package server

import (
	"net/http"
	"time"

	"github.com/SteGG200/storage/internal/handler"
	"github.com/SteGG200/storage/internal/middleware"
)

// Config contains the configuration parameters for the server.
type Config struct {
	Addr        string
	StorageRoot string
}

// NewServer initializes and configures the http.Server.
func NewServer(cfg Config) (*http.Server, error) {
	h, err := handler.NewHandler(cfg.StorageRoot)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	// Register routes with Go 1.22+ method routing
	mux.HandleFunc("GET /src/{path...}", h.GetSrc)
	mux.HandleFunc("POST /src/{path...}", h.PostSrc)
	mux.HandleFunc("PUT /src/{path...}", h.PutSrc)
	mux.HandleFunc("DELETE /src/{path...}", h.DeleteSrc)
	mux.HandleFunc("POST /upload/{path...}", h.PrepareUpload)
	mux.HandleFunc("PATCH /upload/{path...}", h.StreamUpload)
	mux.HandleFunc("GET /progress/{path...}", h.UploadProgress)
	mux.HandleFunc("GET /download/{path...}", h.DownloadFile)

	// Chain middlewares
	handlerWithMiddleware := middleware.Logger(mux)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  2 * time.Hour, // Large write/read timeout for large file uploads
		WriteTimeout: 2 * time.Hour,
		IdleTimeout:  120 * time.Second,
	}

	return srv, nil
}
