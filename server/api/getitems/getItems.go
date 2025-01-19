package getitems

import (
	"encoding/json"
	"net/http"

	"github.com/SteGG200/storage/config"
)

type ListHandler struct {
	config config.Config
}

func New(router *http.ServeMux, config config.Config) {
	listHandler := &ListHandler{
		config: config,
	}

	router.Handle("/", listHandler)
}

func (listHandler *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(listHandler.config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
