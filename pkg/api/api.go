package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"news/pkg/storage"
)

type API struct {
	DB     storage.Storage
	Router *mux.Router
}

func New(db storage.Storage) *API {
	api := API{
		DB:     db,
		Router: mux.NewRouter(),
	}
	api.endpoints()

	return &api
}

func (api *API) endpoints() {
	api.Router.Use(api.headerMiddleware)
	api.Router.HandleFunc("/news/{n}", api.postsHandler).Methods(http.MethodGet)
}

// postsHandler handles GET requests to /news/{n} and returns n latest posts from
// the underlying storage in JSON format.
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	nStr := mux.Vars(r)["n"]
	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	posts, err := api.DB.Posts(n)
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
