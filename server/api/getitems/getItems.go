package getitems

import (
	"encoding/json"
	"net/http"

	"github.com/SteGG200/storage/server/config"
	"github.com/SteGG200/storage/server/mux"
)

type Mux struct {
	*mux.Mux
}

func New(config config.Config) (router *Mux) {
	router = &Mux{
		mux.New(config),
	}

	return
}

func (router *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(router.Config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
