package getitems

import (
	"encoding/json"
	"fmt"
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

	router.HandleFunc("/{path...}", router.ServeData)

	return
}

func (router *Mux) ServeData(w http.ResponseWriter, r *http.Request) {
	items, err := listItems(fmt.Sprintf("%s/%s", router.Config.StoragePath, r.PathValue("path")))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(items)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
