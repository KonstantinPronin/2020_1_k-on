package main

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/server"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/infrastructure"
	"github.com/labstack/echo"
	"log"
)

const port = ":8080"

func main() {
	logger, err := infrastructure.InitLog()
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

	e := echo.New()
	db, err := infrastructure.InitDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf(`error '%s' while closing resource`, err)
		}
	}()

	srv := server.NewServer(port, e, db, rd, logger)
	log.Fatal(srv.ListenAndServe())
}
