package server

import (
	"net/http"

	"github.com/SteGG200/storage/server/api/getitems"
	"github.com/SteGG200/storage/server/config"
)

func NewRouter(config config.Config) (router *http.ServeMux) {
	router = http.NewServeMux()

	// Upload
	// Download
	// Delete

	// List
	router.Handle("/", getitems.New(config))
	// Search
	// Authentication

	return
}
