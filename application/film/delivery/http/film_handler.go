package http

import (
	"2020_1_k-on/application/film"
	"2020_1_k-on/application/models"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"strconv"
)

type FilmHandler struct {
	usecase film.Usecase
}

func NewFilmHandler(e *echo.Echo, usecase film.Usecase) {
	handler := &FilmHandler{
		usecase: usecase,
	}
	e.GET("/films/:id", handler.GetFilm)
	e.GET("/films", handler.GetFilmList)
	e.POST("/films", handler.CreateFilm)
}

func (fh FilmHandler) GetFilm(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp, _ := models.Generate(400, "not number").MarshalJSON()
		ctx.Response().Status = 400
		ctx.Response().Write(resp)
		return err
	}
	f, ok := fh.usecase.GetFilm(uint(id))
	if !ok {
		resp, _ := models.Generate(404, "Not Found").MarshalJSON()
		ctx.Response().Status = 404
		ctx.Response().Write(resp)
		return err
	}
	resp, _ := models.Generate(200, f).MarshalJSON()
	ctx.Response().Status = 200
	_, err = ctx.Response().Write(resp)
	return err
}

func (fh FilmHandler) GetFilmList(ctx echo.Context) error {
	f := fh.usecase.GetFilmsList()
	resp, _ := models.Generate(200, f).MarshalJSON()
	ctx.Response().Status = 200
	_, err := ctx.Response().Write(resp)
	return err
}

func (fh FilmHandler) CreateFilm(ctx echo.Context) error {
	film := new(models.Film)
	defer ctx.Request().Body.Close()
	err := easyjson.UnmarshalFromReader(ctx.Request().Body, film)
	if err != nil {
		resp, _ := models.Generate(400, "request parser error").MarshalJSON()
		ctx.Response().Status = 400
		_, err = ctx.Response().Write(resp)
		return err
	}
	f, ok := fh.usecase.CreateFilm(*film)
	if !ok {
		resp, _ := models.Generate(500, "can't create film").MarshalJSON()
		ctx.Response().Status = 500
		_, err = ctx.Response().Write(resp)
		return err
	}
	resp, _ := models.Generate(200, f).MarshalJSON()
	ctx.Response().Status = 200
	_, err = ctx.Response().Write(resp)
	return err
}
