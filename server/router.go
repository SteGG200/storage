package server

import (
	"net/http"

	"github.com/SteGG200/storage/server/api/download"
	"github.com/SteGG200/storage/server/api/getitems"
	"github.com/SteGG200/storage/server/api/upload"
	"github.com/SteGG200/storage/server/config"
)

func NewRouter(config *config.Config) (router *http.ServeMux) {
	router = http.NewServeMux()

	// Upload
	router.Handle("POST /upload/", http.StripPrefix("/upload", upload.New(config)))
	// Download
	router.Handle("GET /download/", http.StripPrefix("/download", download.New(config)))
	// Delete

	// List
	router.Handle("GET /get/", http.StripPrefix("/get", getitems.New(config)))
	// Search
	// Authentication

	return
}
