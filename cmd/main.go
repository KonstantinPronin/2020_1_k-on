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
	defer logger.Sync()

	e := echo.New()
	db, err := infrastructure.InitGorm()
	if err != nil {
		log.Fatal(err.Error())
	}
	server.NewServer(e, db, logger)
	e.Logger.Fatal(e.Start(":8080"))

}
