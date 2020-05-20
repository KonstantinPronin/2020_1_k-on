package main

import (
	auth "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/server"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/infrastructure"
	"log"
)

const (
	Port1 = ":8081"
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

	rd, err := infrastructure.InitRedis()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := rd.Close(); err != nil {
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

	srv1 := auth.NewServer(Port1, db, rd, logger)
	if err = srv1.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
