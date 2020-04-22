package server

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewAuthServer(usecase user.UseCase) *AuthServer {
	return &AuthServer{usecase: usecase}
}

type AuthServer struct {
	usecase user.UseCase
}

func (a *AuthServer) Login(ctx context.Context, usr *api.User) (*api.Session, error) {
	sessionId, csrfToken, err := a.usecase.Login(usr.Login, usr.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &api.Session{
		SessionId: sessionId,
		CsrfToken: csrfToken,
	}, nil
}

func (a *AuthServer) Check(ctx context.Context, s *api.SessionId) (*api.UserId, error) {
	uid, err := a.usecase.Check(s.SessionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &api.UserId{UserId: uint64(uid)}, nil
}

func (a *AuthServer) Logout(ctx context.Context, sid *api.SessionId) (*api.Empty, error) {
	err := a.usecase.Logout(sid.SessionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.Empty{}, nil
}
