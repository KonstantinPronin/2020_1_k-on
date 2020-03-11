package main

import (
	"github.com/go-park-mail-ru/2020_1_k-on/internal/server"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/infrastructure"
	"github.com/labstack/echo"
	"log"
)

func main() {
	e := echo.New()
	cp, err := infrastructure.InitDatabaseConnection()
	if err != nil {
		log.Fatal(err.Error())
	}

	server.NewServer(e, cp)

	e.Logger.Fatal(e.Start(":8080"))
}
