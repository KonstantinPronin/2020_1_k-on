package server

import (
	filmHandler "github.com/go-park-mail-ru/2020_1_k-on/application/film/delivery/http"
	filmRepository "github.com/go-park-mail-ru/2020_1_k-on/application/film/repository"
	filmUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/film/usecase"
	personHandler "github.com/go-park-mail-ru/2020_1_k-on/application/person/delivery/http"
	personRepository "github.com/go-park-mail-ru/2020_1_k-on/application/person/repository"
	personUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/person/usecase"
	reviewHandler "github.com/go-park-mail-ru/2020_1_k-on/application/review/delivery/http"
	reviewRepository "github.com/go-park-mail-ru/2020_1_k-on/application/review/repository"
	reviewUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/review/usecase"
	serialHandler "github.com/go-park-mail-ru/2020_1_k-on/application/series/delivery/http"
	serialRepository "github.com/go-park-mail-ru/2020_1_k-on/application/series/repository"
	serialUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/series/usecase"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	session "github.com/go-park-mail-ru/2020_1_k-on/application/session/repository"
	userHandler "github.com/go-park-mail-ru/2020_1_k-on/application/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2020_1_k-on/application/user/repository"
	userUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/user/usecase"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	middleware2 "github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

type Server struct {
	port string
	e    *echo.Echo
}

func NewServer(port string, e *echo.Echo, db *gorm.DB, rd *redis.Client, logger *zap.Logger) *Server {
	//middleware
	e.Use(middleware.Middleware)
	//e.Use(middleware.CORS)
	e.Use(middleware2.CORSWithConfig(middleware2.CORSConfig{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	//user handler
	sessions := session.NewSessionDatabase(rd, logger)
	users := userRepository.NewUserDatabase(db, logger)
	auth := middleware.NewAuth(sessions)
	user := userUsecase.NewUser(sessions, users, logger)
	userHandler.NewUserHandler(e, user, auth, logger)

	//person handler
	persons := personRepository.NewPersonDatabase(db, logger)
	person := personUsecase.NewPerson(persons, logger)
	personHandler.NewPersonHandler(e, person, logger)

	//series handler
	series := serialRepository.NewPostgresForSeries(db)
	seriesUC := serialUsecase.NewSeriesUsecase(series)
	serialHandler.NewSeriesHandler(e, seriesUC, person)

	//film handler
	films := filmRepository.NewPostgresForFilms(db)
	film := filmUsecase.NewFilmUsecase(films)
	filmHandler.NewFilmHandler(e, film, person)

	//review handler
	filmReviewsRep := reviewRepository.NewFilmReviewDatabase(db, logger)
	filmReview := reviewUsecase.NewFilmReview(filmReviewsRep, films)
	seriesReviewsRep := reviewRepository.NewSeriesReviewDatabase(db, logger)
	seriesReview := reviewUsecase.NewSeriesReview(seriesReviewsRep, series)
	reviewHandler.NewReviewHandler(e, filmReview, seriesReview, auth, logger)

	return &Server{
		port: port,
		e:    e,
	}
}

func (s Server) ListenAndServe() error {
	return s.e.Start(s.port)
}
