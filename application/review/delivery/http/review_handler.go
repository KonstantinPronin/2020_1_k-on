package http

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	"github.com/go-park-mail-ru/2020_1_k-on/application/session"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ReviewHandler struct {
	uc     review.UseCase
	logger *zap.Logger
}

func NewReviewHandler(e *echo.Echo, uc review.UseCase, auth middleware.Auth, logger *zap.Logger) {
	handler := ReviewHandler{
		uc:     uc,
		logger: logger,
	}

	e.POST("/review", handler.Add, auth.GetSession, middleware.ParseErrors)
	e.GET("/films/:id/reviews", handler.GetByFilm, middleware.ParseErrors)
}

func (r *ReviewHandler) Add(ctx echo.Context) error {
	rev := new(models.Review)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, rev); err != nil {
		r.logger.Error("request parser error")
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "request parser error")
	}

	rev.UserId = ctx.Get(session.UserIdKey).(uint)
	if err := r.uc.Add(rev); err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, nil)
}

func (r *ReviewHandler) GetByFilm(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 0 {
		r.logger.Error("wrong parameter")
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "wrong parameter")
	}

	reviews, err := r.uc.GetByFilmId(uint(id))
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, reviews)
}
