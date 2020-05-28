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
	client3 "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/client"
	personHandler "github.com/go-park-mail-ru/2020_1_k-on/application/person/delivery/http"
	personRepository "github.com/go-park-mail-ru/2020_1_k-on/application/person/repository"
	personUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/person/usecase"
	playlistHandler "github.com/go-park-mail-ru/2020_1_k-on/application/playlist/delivery/http"
	playlistRepository "github.com/go-park-mail-ru/2020_1_k-on/application/playlist/repository"
	playlistUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/playlist/usecase"
	reviewHandler "github.com/go-park-mail-ru/2020_1_k-on/application/review/delivery/http"
	reviewRepository "github.com/go-park-mail-ru/2020_1_k-on/application/review/repository"
	reviewUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/review/usecase"
	serialHandler "github.com/go-park-mail-ru/2020_1_k-on/application/series/delivery/http"
	serialRepository "github.com/go-park-mail-ru/2020_1_k-on/application/series/repository"
	serialUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/series/usecase"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	subsHandler "github.com/go-park-mail-ru/2020_1_k-on/application/subscription/delivery/http"
	subsRepository "github.com/go-park-mail-ru/2020_1_k-on/application/subscription/repository"
	subsUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/subscription/usecase"
	userHandler "github.com/go-park-mail-ru/2020_1_k-on/application/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2020_1_k-on/application/user/repository"
	userUsecase "github.com/go-park-mail-ru/2020_1_k-on/application/user/usecase"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/conf"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/constants"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	middleware2 "github.com/labstack/echo/middleware"
	"github.com/microcosm-cc/bluemonday"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"log"
	_ "net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	rpcAuth         *client.AuthClient
	rpcFilmFilter   *client2.FilmFilterClient
	rpcSeriesFilter *client3.SeriesFilterClient
	port            string
	e               *echo.Echo
}

func NewServer(srvConf *conf.Service, e *echo.Echo, db *gorm.DB, logger *zap.Logger) *Server {
	//tracing
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "main_server",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, _, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		log.Fatal("cannot create tracer", err)
	}

	opentracing.SetGlobalTracer(tracer)

	//microservices
	rpcSeriesFilter, err := client3.NewSeriesFilterClient(srvConf.Host, srvConf.Port3, logger, tracer)
	if err != nil {
		log.Fatal(err.Error())
	}
	rpcFilmFilter, err := client2.NewFilmFilterClient(srvConf.Host, srvConf.Port2, logger, tracer)
	if err != nil {
		log.Fatal(err.Error())
	}
	rpcAuth, err := client.NewAuthClient(srvConf.Host, srvConf.Port1, logger, tracer)
	if err != nil {
		log.Fatal(err.Error())
	}

	//middleware
	sanitizer := bluemonday.UGCPolicy()
	ioLog := middleware.NewLogger(logger)

	e.Use(ioLog.Log)
	//e.Use(middleware.CORS)
	e.Use(middleware2.CORSWithConfig(middleware2.CORSConfig{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType,
			echo.HeaderAccept, constants.CSRFHeader},
		AllowCredentials: true,
	}))

	//user handler
	users := userRepository.NewUserDatabase(db, logger)
	auth := middleware.NewAuth(rpcAuth)
	user := userUsecase.NewUser(users, logger)
	err = userHandler.NewUserHandler(e, rpcAuth, user, auth, logger, sanitizer)
	if err != nil {
		log.Fatal(err.Error())
	}

	//person handler
	persons := personRepository.NewPersonDatabase(db, logger)
	person := personUsecase.NewPerson(persons, logger)
	personHandler.NewPersonHandler(e, person, logger, sanitizer)

	//series handler
	series := serialRepository.NewPostgresForSeries(db)
	seriesUC := serialUsecase.NewSeriesUsecase(series)
	serialHandler.NewSeriesHandler(e, rpcSeriesFilter, seriesUC, person)

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

	//playlist handler
	playlists := playlistRepository.NewPlaylistDatabase(db, logger)
	playlist := playlistUsecase.NewPlaylist(playlists, logger)
	playlistHandler.NewPlaylistHandler(e, playlist, auth, logger, sanitizer)

	//subscription handler
	subsRep := subsRepository.NewSubscriptionDatabase(db, logger)
	subs := subsUsecase.NewSubscription(playlists, subsRep, logger)
	subsHandler.NewSubscriptionHandler(e, subs, auth, logger)

	//prometeus

	prometheus.MustRegister(middleware.FooCount, middleware.Hits)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return &Server{
		rpcAuth:         rpcAuth,
		rpcFilmFilter:   rpcFilmFilter,
		rpcSeriesFilter: rpcSeriesFilter,
		port:            srvConf.Port0,
		e:               e,
	}
}

func (s Server) ListenAndServe() error {
	defer func() {
		s.rpcAuth.Close()
		s.rpcFilmFilter.Close()
		s.rpcSeriesFilter.Close()
	}()
	return s.e.Start(s.port)
}
