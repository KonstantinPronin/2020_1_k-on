package http

import (
	"errors"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/serial"
	"github.com/labstack/echo"
	"strconv"
)

type SerialHandler struct {
	usecase serial.Usecase
}

func NewSerialHandler(e *echo.Echo, usecase serial.Usecase) {
	handler := &SerialHandler{
		usecase: usecase,
	}
	e.GET("/serial/:id", handler.GetSerial)
	e.GET("/serial/:id/seasons", handler.GetSerialSeasons)
	e.GET("/seasons/:id/series", handler.GetSerialSeries)
}

func (sh SerialHandler) GetSerial(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp, _ := models.Generate(400, "not number", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return err
	}
	serial, ok := sh.usecase.GetSerialByID(uint(id))
	if !ok {
		resp, _ := models.Generate(404, "Not Found", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("Not Found")
	}
	resp, _ := models.Generate(200, serial, &ctx).MarshalJSON()
	_, err = ctx.Response().Write(resp)
	return err
}

func (sh SerialHandler) GetSerialSeasons(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp, _ := models.Generate(400, "not number", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return err
	}
	seasons, ok := sh.usecase.GetSerialSeasons(uint(id))
	if !ok {
		resp, _ := models.Generate(404, "Not Found", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("Not Found")
	}
	resp, _ := models.Generate(200, seasons, &ctx).MarshalJSON()
	_, err = ctx.Response().Write(resp)
	return err
}

func (sh SerialHandler) GetSerialSeries(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp, _ := models.Generate(400, "not number", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return err
	}
	seasons, ok := sh.usecase.GetSeasonSeries(uint(id))
	if !ok {
		resp, _ := models.Generate(404, "Not Found", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("Not Found")
	}
	resp, _ := models.Generate(200, seasons, &ctx).MarshalJSON()
	_, err = ctx.Response().Write(resp)
	return err
}
