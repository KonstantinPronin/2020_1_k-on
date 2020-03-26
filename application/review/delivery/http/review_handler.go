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
	film   review.UseCase
	serial review.UseCase
	logger *zap.Logger
}

func NewReviewHandler(e *echo.Echo, review review.UseCase, serial review.UseCase, auth middleware.Auth, logger *zap.Logger) {
	handler := ReviewHandler{
		film:   review,
		serial: serial,
		logger: logger,
	}

	e.POST("/films/review", handler.AddFilmReview, auth.GetSession, middleware.ParseErrors)
	e.POST("/serials/review", handler.AddSerialReview, auth.GetSession, middleware.ParseErrors)
	e.GET("/films/:id/reviews", handler.GetByFilm, middleware.ParseErrors)
	e.GET("/serials/:id/reviews", handler.GetBySerial, middleware.ParseErrors)
}

func (r *ReviewHandler) AddFilmReview(ctx echo.Context) error {
	rev, err := r.parseRequestBody(ctx)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "request parser error")
	}

	if err := r.film.Add(rev); err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, nil)
}

func (r *ReviewHandler) AddSerialReview(ctx echo.Context) error {
	rev, err := r.parseRequestBody(ctx)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "request parser error")
	}

	if err := r.serial.Add(rev); err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, nil)
}

func (r *ReviewHandler) GetByFilm(ctx echo.Context) error {
	id, err := r.getProductId(ctx)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "wrong parameter")
	}

	reviews, err := r.film.GetByProductId(uint(id))
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, reviews)
}

func (r *ReviewHandler) GetBySerial(ctx echo.Context) error {
	id, err := r.getProductId(ctx)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "wrong parameter")
	}

	reviews, err := r.serial.GetByProductId(uint(id))
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, reviews)
}

func (r *ReviewHandler) parseRequestBody(ctx echo.Context) (*models.Review, error) {
	rev := new(models.Review)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, rev); err != nil {
		r.logger.Error("request parser error")
		return nil, err
	}

	rev.UserId = ctx.Get(session.UserIdKey).(uint)
	return rev, nil
}

func (r *ReviewHandler) getProductId(ctx echo.Context) (uint, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 0 {
		r.logger.Error("wrong parameter")
		return 0, err
	}

	return uint(id), nil
}
