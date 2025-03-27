package main

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"

	"news/pkg/api"
	"news/pkg/storage"
	"news/pkg/storage/memdb"
)

type server struct {
	api *api.API
}

func main() {
	var (
		srv server
		db  storage.Storage

		dev bool
	)
	flag.BoolVar(&dev, "dev", false, "Run the server in development mode with in-memory DB.")
	flag.Parse()

	switch dev {
	case false:
		log.Error("Not implemented")
		return
	case true:
		log.Info("Run server with in memory DB")
		db = memdb.New()

		// TODO: remove after RSS implementation.
		posts, err := memdb.LoadTestPosts("test_data/post_examples.json")
		if err != nil {
			log.Fatal(err)
		}
		n, err := db.AddPosts(posts)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Posts in DB: %d", n)
	}

	srv.api = api.New(db)
	http.ListenAndServe(":8080", srv.api.Router())
}
