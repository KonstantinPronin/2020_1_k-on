package server

import (
	filmHandler "github.com/go-park-mail-ru/2020_1_k-on/application/film/delivery/http"
	filmRepository "github.com/go-park-mail-ru/2020_1_k-on/application/film/repository"
	filmUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/film/usecase"
	imageHandler "github.com/go-park-mail-ru/2020_1_k-on/application/image/delivery/http"
	imageRepository "github.com/go-park-mail-ru/2020_1_k-on/application/image/repository"
	imageUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/image/usecase"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/client"
	client2 "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/client"
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
	userHandler "github.com/go-park-mail-ru/2020_1_k-on/application/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2020_1_k-on/application/user/repository"
	userUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/user/usecase"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/conf"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	middleware2 "github.com/labstack/echo/middleware"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"log"
)

type Server struct {
	rpcAuth       *client.AuthClient
	rpcFilmFilter *client2.FilmFilterClient
	port          string
	e             *echo.Echo
}

func NewServer(srvConf *conf.Service, e *echo.Echo, db *gorm.DB, logger *zap.Logger) *Server {
	//microservices
	rpcAuth, err := client.NewAuthClient(srvConf.Host, srvConf.Port1, logger)
	if err != nil {
		log.Fatal(err.Error())
	}
	rpcFilmFilter, err := client2.NewFilmFilterClient(srvConf.Host, srvConf.Port2, logger)
	if err != nil {
		log.Fatal(err.Error())
	}

	//middleware
	sanitizer := bluemonday.UGCPolicy()
	ioLog := middleware.NewLogger(logger)

	e.Use(ioLog.Log)
	//e.Use(middleware.CORS)
	e.Use(middleware2.CORSWithConfig(middleware2.CORSConfig{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken},
		AllowCredentials: true,
	}))

	//user handler
	users := userRepository.NewUserDatabase(db, logger)
	auth := middleware.NewAuth(rpcAuth)
	user := userUsecase.NewUser(users, logger)
	userHandler.NewUserHandler(e, rpcAuth, user, auth, logger, sanitizer)

	//person handler
	persons := personRepository.NewPersonDatabase(db, logger)
	person := personUsecase.NewPerson(persons, logger)
	personHandler.NewPersonHandler(e, person, logger, sanitizer)

	//series handler
	series := serialRepository.NewPostgresForSeries(db)
	seriesUC := serialUsecase.NewSeriesUsecase(series)
	serialHandler.NewSeriesHandler(e, seriesUC, person)

	//film handler
	films := filmRepository.NewPostgresForFilms(db)
	film := filmUsecase.NewFilmUsecase(films)
	filmHandler.NewFilmHandler(e, rpcFilmFilter, film, person, sanitizer)

	//review handler
	filmReviewsRep := reviewRepository.NewFilmReviewDatabase(db, logger)
	filmReview := reviewUsecase.NewFilmReview(filmReviewsRep, films)
	seriesReviewsRep := reviewRepository.NewSeriesReviewDatabase(db, logger)
	seriesReview := reviewUsecase.NewSeriesReview(seriesReviewsRep, series)
	reviewHandler.NewReviewHandler(e, filmReview, seriesReview, auth, logger, sanitizer)

	//image handler
	images := imageRepository.NewImageRepository()
	image := imageUsecase.NewImage(images, logger)
	imageHandler.NewUserHandler(e, image, user, auth, logger)

	return &Server{
		rpcAuth:       rpcAuth,
		rpcFilmFilter: rpcFilmFilter,
		port:          srvConf.Port0,
		e:             e,
	}
}

func (s Server) ListenAndServe() error {
	defer func() {
		s.rpcAuth.Close()
		s.rpcFilmFilter.Close()
	}()
	return s.e.Start(s.port)
}
