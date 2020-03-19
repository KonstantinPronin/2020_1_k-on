package server

import (
	filmHandler "github.com/go-park-mail-ru/2020_1_k-on/application/film/delivery/http"
	filmRepository "github.com/go-park-mail-ru/2020_1_k-on/application/film/repository"
	filmUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/film/usecase"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	session "github.com/go-park-mail-ru/2020_1_k-on/application/session/repository"
	userHandler "github.com/go-park-mail-ru/2020_1_k-on/application/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2020_1_k-on/application/user/repository"
	userUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/user/usecase"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Server struct {
	port string
	e    *echo.Echo
}

func NewServer(port string, e *echo.Echo, db *gorm.DB, rd *redis.Client, logger *zap.Logger) *Server {
	//middleware
	e.Use(middleware.Middleware)
	e.Use(middleware.CORS)

	//film handler
	films := filmRepository.NewPostgresForFilms(db)
	film := filmUsecase.NewFilmUsecase(films)
	filmHandler.NewFilmHandler(e, film)

	//user handler
	sessions := session.NewSessionDatabase(rd, logger)
	users := userRepository.NewUserDatabase(db, logger)
	auth := middleware.NewAuth(sessions)
	user := userUsecase.NewUser(sessions, users, logger)
	userHandler.NewUserHandler(e, user, auth, logger)

	return &Server{
		port: port,
		e:    e,
	}
}

func (s Server) ListenAndServe() error {
	return s.e.Start(s.port)
}
