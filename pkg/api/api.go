package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"news/pkg/storage"
)

type API struct {
	db storage.Storage
	r  *mux.Router
}

func New(db storage.Storage) *API {
	api := API{
		db: db,
		r:  mux.NewRouter(),
	}
	api.endpoints()

	return &api
}

func (api *API) endpoints() {
	api.r.Use(api.headerMiddleware)
	api.r.HandleFunc("/news/{n}", api.postsHandler).Methods(http.MethodGet)
}

func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	nStr := mux.Vars(r)["n"]
	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	posts, err := api.db.Posts(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
