package server

import (
	api "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/filter/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/filter/usecase"
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
	filter *SeriesFilter
}

func NewServer(port string, db *gorm.DB, logger *zap.Logger) *Server {
	rep := repository.NewSeriesFiltersDb(db)
	use := usecase.NewSeriesFilter(rep)

	return &Server{
		port:   port,
		filter: NewSeriesFilter(use),
	}
}

func (s *Server) ListenAndServe() error {
	//tracing
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "series_server",
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
	//defer closer.Close()
	//
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	gServer := grpc.NewServer(grpc.
		UnaryInterceptor(traceutils.OpenTracingServerInterceptor(tracer)))
	api.RegisterSeriesFilterServer(gServer, s.filter)

	err = gServer.Serve(listener)
	if err != nil {
		return nil
	}

	return nil
}
