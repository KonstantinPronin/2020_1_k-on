package server

import (
	api "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/filter/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/filter/usecase"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
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
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	gServer := grpc.NewServer()
	api.RegisterSeriesFilterServer(gServer, s.filter)

	err = gServer.Serve(listener)
	if err != nil {
		return nil
	}

	return nil
}
