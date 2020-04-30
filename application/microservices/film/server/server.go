package server

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/filter/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/filter/usecase"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"

	traceutils "github.com/opentracing-contrib/go-grpc"
)

type Server struct {
	port   string
	filter *FilmFilter
}

func NewServer(port string, db *gorm.DB, logger *zap.Logger) *Server {
	films := repository.NewFilmFiltersDb(db)
	film := usecase.NewFilmFilter(films)

	return &Server{
		port:   port,
		filter: NewFilmFilter(film),
	}
}

func (s *Server) ListenAndServe() error {
	//tracing
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "film_server",
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
	//

	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	gServer := grpc.NewServer(grpc.
		UnaryInterceptor(traceutils.OpenTracingServerInterceptor(tracer)))
	api.RegisterFilmFilterServer(gServer, s.filter)

	err = gServer.Serve(listener)
	if err != nil {
		return nil
	}

	return nil
}
