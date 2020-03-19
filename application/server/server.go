package server

import (
	"2020_1_k-on/application/film/delivery/http"
	"2020_1_k-on/application/film/repository"
	"2020_1_k-on/application/film/usecase"
	"github.com/go-park-mail-ru/2020_1_k-on/application/middleware"
	repository2 "github.com/go-park-mail-ru/2020_1_k-on/application/session/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user/delivery/http"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user/usecase"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Server struct {
	port string
	e    *echo.Echo
}

func NewServer(port string, db *gorm.DB, rd *redis.Client, logger *zap.Logger) *server {
	e := echo.New()

	//middleware
	e.Use(Middleware)
	e.Use(CORS)

	//film handler
	filmrepo := repository.NewPostgresForFilms(db)
	filmUsecase := usecase.NewFilmUsecase(filmrepo)
	http.NewFilmHandler(router, filmUsecase)

	//user handler
	sessions := repository2.NewSessionDatabase(rd, logger)
	users := repository.NewUserDatabase(db, logger)
	auth := middleware.NewAuth(sessions)
	user := usecase.NewUser(sessions, users, logger)
	http.NewUserHandler(e, user, auth, logger)

	return &Server{
		port: port,
		e:    e,
	}
}

func (s Server) ListenAndServe() error {
	return http.ListenAndServe(s.port, s.e)
}
