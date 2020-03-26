package http

import (
	"errors"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/labstack/echo"
	"strconv"
)

type SeriesHandler struct {
	usecase series.Usecase
}

func NewSeriesHandler(e *echo.Echo, usecase series.Usecase) {
	handler := &SeriesHandler{
		usecase: usecase,
	}
	e.GET("/series/:id", handler.GetSeries)
	e.GET("/series/:id/seasons", handler.GetSeriesSeasons)
	e.GET("/seasons/:id/series", handler.GetSeasonEpisodes)
}

func (sh SeriesHandler) GetSeries(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp, _ := models.Generate(400, "not number", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return err
	}
	serial, ok := sh.usecase.GetSeriesByID(uint(id))
	if !ok {
		resp, _ := models.Generate(404, "Not Found", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("Not Found")
	}
	resp, _ := models.Generate(200, serial, &ctx).MarshalJSON()
	_, err = ctx.Response().Write(resp)
	return err
}

func (sh SeriesHandler) GetSeriesSeasons(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp, _ := models.Generate(400, "not number", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return err
	}
	seasons, ok := sh.usecase.GetSeriesSeasons(uint(id))
	if !ok {
		resp, _ := models.Generate(404, "Not Found", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("Not Found")
	}
	resp, _ := models.Generate(200, seasons, &ctx).MarshalJSON()
	_, err = ctx.Response().Write(resp)
	return err
}

func (sh SeriesHandler) GetSeasonEpisodes(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp, _ := models.Generate(400, "not number", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return err
	}
	seasons, ok := sh.usecase.GetSeasonEpisodes(uint(id))
	if !ok {
		resp, _ := models.Generate(404, "Not Found", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("Not Found")
	}
	resp, _ := models.Generate(200, seasons, &ctx).MarshalJSON()
	_, err = ctx.Response().Write(resp)
	return err
}
