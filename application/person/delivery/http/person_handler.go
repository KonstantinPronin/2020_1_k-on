package http

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/person"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PersonHandler struct {
	usecase person.UseCase
	logger  *zap.Logger
}

func NewPersonHandler(e *echo.Echo, usecase person.UseCase, logger *zap.Logger) {
	handler := PersonHandler{
		usecase: usecase,
		logger:  logger,
	}

	e.GET("/persons/:id", handler.GetById, middleware.ParseErrors)
	e.POST("/persons", handler.Add, middleware.ParseErrors)
	e.PUT("/persons", handler.Update, middleware.ParseErrors)
}

func (handler *PersonHandler) GetById(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 0 {
		handler.logger.Error("wrong parameter")
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "wrong parameter")
	}

	p, err := handler.usecase.GetById(uint(id))
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, p)
}

func (handler *PersonHandler) Add(ctx echo.Context) error {
	p := new(models.Person)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, p); err != nil {
		handler.logger.Error("request parser error")
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "request parser error")
	}

	p, err := handler.usecase.Add(p)
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, p)
}

func (handler *PersonHandler) Update(ctx echo.Context) error {
	p := new(models.Person)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, p); err != nil {
		handler.logger.Error("request parser error")
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "request parser error")
	}

	p, err := handler.usecase.Update(p)
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, p)
}
