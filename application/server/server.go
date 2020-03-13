package server

import (
	http1 "2020_1_k-on/application/film/delivery/http"
	"2020_1_k-on/application/film/repository"
	"2020_1_k-on/application/film/usecase"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
)

type server struct {
	port   string
	router *echo.Echo
}

func NewServer(port string, connection *gorm.DB) *server {
	router := echo.New()
	router.Use(Middleware)
	filmrepo := repository.NewPostgresForFilms(connection)
	filmUsecase := usecase.NewUserUsecase(filmrepo)

	http1.NewFilmHandler(router, filmUsecase)

	return &server{
		port:   port,
		router: router,
	}
}

func (serv server) ListenAndServe() error {
	return http.ListenAndServe(serv.port, serv.router)
}
