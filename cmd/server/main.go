package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"news/pkg/api"
	"news/pkg/rss"
	"news/pkg/storage"
	"news/pkg/storage/memdb"
	"news/pkg/storage/postgres"
)

func main() {
	var (
		sdb storage.Storage
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
		conf := postgres.Config{
			User:     "postgres",
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			DBName:   "news",
		}
		db, err := postgres.New(conf.ConString())
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			log.Fatal(fmt.Errorf("%w: %v", storage.ErrDBNotResponding, err))
		}
		log.Infof("connected to postgres: %s", conf)
		sdb = db

	case true:
		log.Info("Run server with in memory DB")
		sdb = memdb.New()
	}

	conf, err := rss.LoadConf("cmd/server/config.json")
	if err != nil {
		log.Fatalf("unable to load RSS parser config: %v", err)
	}

	api := api.New(sdb)
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
				err := api.DB.AddPosts(msg.Data)
				if err != nil {
					log.Warnf("Error while adding posts from %s to DB: %v", msg.Source, err)
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
	http.ListenAndServe(":8088", api.Router)
	wg.Wait()
}
