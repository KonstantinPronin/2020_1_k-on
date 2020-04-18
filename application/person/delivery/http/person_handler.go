package http

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/person"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PersonHandler struct {
	usecase   person.UseCase
	logger    *zap.Logger
	sanitizer *bluemonday.Policy
}

func NewPersonHandler(e *echo.Echo, usecase person.UseCase, logger *zap.Logger, sanitizer *bluemonday.Policy) {
	handler := PersonHandler{
		usecase:   usecase,
		logger:    logger,
		sanitizer: sanitizer,
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

	handler.sanitize(p)
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

	handler.sanitize(p)
	p, err := handler.usecase.Update(p)
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, p)
}

func (handler *PersonHandler) sanitize(p *models.Person) {
	p.Name = handler.sanitizer.Sanitize(p.Name)
	p.Image = handler.sanitizer.Sanitize(p.Image)
	p.BirthPlace = handler.sanitizer.Sanitize(p.BirthPlace)
	p.BirthDate = handler.sanitizer.Sanitize(p.BirthDate)
	p.Occupation = handler.sanitizer.Sanitize(p.Occupation)
}
