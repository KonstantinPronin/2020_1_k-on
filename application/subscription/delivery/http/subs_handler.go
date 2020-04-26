package http

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	"github.com/go-park-mail-ru/2020_1_k-on/application/session"
	"github.com/go-park-mail-ru/2020_1_k-on/application/subscription"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type SubscriptionHandler struct {
	useCase subscription.UseCase
	logger  *zap.Logger
}

func NewSubscriptionHandler(e *echo.Echo, useCase subscription.UseCase, auth middleware.Auth, logger *zap.Logger) {
	handler := SubscriptionHandler{
		useCase: useCase,
		logger:  logger,
	}

	e.POST("/subscribe/:pid", handler.Subscribe, auth.GetSession, middleware.ParseErrors)
	e.DELETE("/unsubscribe/:pid", handler.Unsubscribe, auth.GetSession, middleware.ParseErrors)
	e.GET("/subscriptions", handler.Subscription, auth.GetSession, middleware.ParseErrors)
}

func (s *SubscriptionHandler) Subscribe(ctx echo.Context) error {
	pid, err := s.getPlaylistId(ctx)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "wrong parameter")
	}
	userId := ctx.Get(session.UserIdKey).(uint)

	err = s.useCase.Subscribe(pid, userId)
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, "")
}

func (s *SubscriptionHandler) Unsubscribe(ctx echo.Context) error {
	pid, err := s.getPlaylistId(ctx)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "wrong parameter")
	}
	userId := ctx.Get(session.UserIdKey).(uint)

	err = s.useCase.Unsubscribe(pid, userId)
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, "")
}

func (s *SubscriptionHandler) Subscription(ctx echo.Context) error {
	userId := ctx.Get(session.UserIdKey).(uint)

	plist, err := s.useCase.Subscriptions(userId)
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, plist)
}

func (s *SubscriptionHandler) getPlaylistId(ctx echo.Context) (uint, error) {
	id, err := strconv.Atoi(ctx.Param("pid"))
	if err != nil || id < 0 {
		s.logger.Error("wrong parameter")
		return 0, err
	}

	return uint(id), nil
}
