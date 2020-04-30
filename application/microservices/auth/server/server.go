package server

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/api"
	session "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/session/repository"
	userRepository "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/user/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/user/usecase"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	port string
	auth *AuthServer
}

func NewServer(port string, db *gorm.DB, rd *redis.Client, logger *zap.Logger) *Server {
	sessions := session.NewSessionDatabase(rd, logger)
	users := userRepository.NewUserDatabase(db, logger)
	user := usecase.NewUser(sessions, users, logger)

	return &Server{
		port: port,
		auth: NewAuthServer(user),
	}
}

func (s *Server) ListenAndServe() error {
	//tracing
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "auth_server",
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
	api.RegisterAuthServer(gServer, s.auth)

	err = gServer.Serve(listener)
	if err != nil {
		return nil
	}

	return nil
}
