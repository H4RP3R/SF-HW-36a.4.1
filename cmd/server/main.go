package main

import (
	"flag"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"news/pkg/api"
	"news/pkg/rss"
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

	var (
		done    = make(chan struct{})
		msgChan = make(chan rss.ParserMsg)
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
	}

	conf, err := rss.LoadConf("cmd/server/config.json")
	if err != nil {
		log.Fatalf("unable to load RSS parser config: %v", err)
	}

	srv.api = api.New(db)
	parser := rss.NewParser(*conf)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			log.Infof("Message receiver stopped")
			wg.Done()
		}()

		for msg := range msgChan {
			if msg.Err != nil {
				log.Warnf("Error while parsing %s: %v", msg.Source, msg.Err)
			} else {
				err := srv.api.DB.AddPosts(msg.Data)
				if err != nil {
					log.Warnf("Error while adding posts from %s to DB: %v", msg.Source, msg.Err)
				} else {
					log.Infof("DB updated with posts from %s", msg.Source)
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(parser.Delay)

		defer func() {
			close(msgChan)
			ticker.Stop()
			log.Info("Parser stopped")
			wg.Done()
		}()

		for {
			select {
			case <-done:
				return
			default:
				parser.Run(msgChan)
				<-ticker.C
			}
		}
	}()

	// TODO: graceful shutdown.
	http.ListenAndServe(":8080", srv.api.Router())
	wg.Wait()
}
