package server

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/api"
	session "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/session/repository"
	userRepository "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/user/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/user/usecase"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	gServer := grpc.NewServer()
	api.RegisterAuthServer(gServer, s.auth)

	err = gServer.Serve(listener)
	if err != nil {
		return nil
	}

	return nil
}
