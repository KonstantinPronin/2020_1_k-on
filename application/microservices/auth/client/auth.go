package client

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/api"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AuthClient struct {
	client api.AuthClient
	gConn  *grpc.ClientConn
	logger *zap.Logger
}

func NewAuthClient(host, port string, logger *zap.Logger, tracer opentracing.Tracer) (*AuthClient, error) {
	gConn, err := grpc.Dial(
		host+port,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		return nil, err
	}

	return &AuthClient{client: api.NewAuthClient(gConn), gConn: gConn, logger: logger}, nil
}

func (a *AuthClient) Login(login string, password string) (sessionId string, csrfToken string, err error) {
	usr := &api.User{
		Login:    login,
		Password: password,
	}

	session, err := a.client.Login(context.Background(), usr)
	if err != nil {
		return "", "", err
	}

	return session.SessionId, session.CsrfToken, nil
}

func (a *AuthClient) Check(sessionId string) (userId uint, err error) {
	sid := &api.SessionId{SessionId: sessionId}

	uid, err := a.client.Check(context.Background(), sid)
	if err != nil {
		return 0, err
	}

	return uint(uid.UserId), err
}

func (a *AuthClient) Logout(sessionId string) error {
	sid := &api.SessionId{SessionId: sessionId}

	_, err := a.client.Logout(context.Background(), sid)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthClient) Close() {
	if err := a.gConn.Close(); err != nil {
		a.logger.Error("error while closing grpc connection")
	}
}
