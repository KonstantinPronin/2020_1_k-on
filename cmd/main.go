package main

import (
	"github.com/go-park-mail-ru/2020_1_k-on/internal/server"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/infrastructure"
	"github.com/labstack/echo"
	"log"
)

func main() {
	logger, err := infrastructure.InitLog()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	rd, err := infrastructure.InitRedis()
	if err != nil {
		log.Fatal(err.Error())
	}

	e := echo.New()
	db, err := infrastructure.InitGorm()
	if err != nil {
		log.Fatal(err.Error())
	}

	server.NewServer(e, db, rd, logger)
	e.Logger.Fatal(e.Start(":8080"))
}
