package main

import (
	series "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/server"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/infrastructure"
	"log"
)

const (
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

	srv3 := series.NewServer(Port3, db, logger)
	if err = srv3.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
