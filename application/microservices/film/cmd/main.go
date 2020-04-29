package main

import (
	film "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/server"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/infrastructure"
	"log"
)

const (
	Host  = "127.0.0.1"
	Port0 = ":8080"
	Port1 = ":8081"
	Port2 = ":8082"
	Port3 = ":8083"
)

func main() {
	logger, err := infrastructure.InitLog("conf/log.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf(`error '%s' while closing resource`, err)
		}
	}()

	db, err := infrastructure.InitDatabase("conf/db.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf(`error '%s' while closing resource`, err)
		}
	}()

	srv2 := film.NewServer(Port2, db, logger)
	if err = srv2.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
