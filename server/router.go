package server

import (
	"net/http"

	"github.com/SteGG200/storage/config"
	"github.com/SteGG200/storage/server/api/getitems"
)

func NewRouter(config config.Config) (router *http.ServeMux) {
	router = http.NewServeMux()

	// Upload

	getitems.New(router, config)

	// Download
	// Delete
	// List
	// Search
	// Authentication

	return
}
