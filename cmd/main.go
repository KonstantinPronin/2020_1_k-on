package main

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/server"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/conf"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/infrastructure"
	"github.com/labstack/echo"
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

	e := echo.New()
	db, err := infrastructure.InitDatabase("conf/db.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf(`error '%s' while closing resource`, err)
		}
	}()

	srvConf := &conf.Service{
		Host:  Host,
		Port0: Port0,
		Port1: Port1,
		Port2: Port2,
		Port3: Port3,
	}

	srv0 := server.NewServer(srvConf, e, db, logger)
	log.Fatal(srv0.ListenAndServe())
}
