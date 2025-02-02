package server

import (
	"net/http"

	"github.com/SteGG200/storage/server/api/auth"
	"github.com/SteGG200/storage/server/api/create"
	"github.com/SteGG200/storage/server/api/download"
	"github.com/SteGG200/storage/server/api/get"
	"github.com/SteGG200/storage/server/api/remove"
	"github.com/SteGG200/storage/server/api/upload"
	"github.com/SteGG200/storage/server/config"
)

func NewRouter(config *config.Config) (router *http.ServeMux) {
	router = http.NewServeMux()

	// Upload
	router.Handle("POST /upload/", http.StripPrefix("/upload", upload.New(config)))
	// Download
	router.Handle("GET /download/", http.StripPrefix("/download", download.New(config)))
	//Create
	router.Handle("POST /create/", http.StripPrefix("/create", create.New(config)))
	// Delete
	router.Handle("DELETE /remove/", http.StripPrefix("/remove", remove.New(config)))

	// List
	router.Handle("GET /get/", http.StripPrefix("/get", get.New(config)))
	// Search
	// Authentication
	router.Handle("POST /auth/", http.StripPrefix("/auth", auth.New(config)))

	return
}
