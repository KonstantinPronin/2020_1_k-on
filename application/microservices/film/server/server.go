package server

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/filter/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/filter/usecase"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
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
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	gServer := grpc.NewServer()
	api.RegisterFilmFilterServer(gServer, s.filter)

	err = gServer.Serve(listener)
	if err != nil {
		return nil
	}

	return nil
}
