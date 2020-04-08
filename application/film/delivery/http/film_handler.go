package http

import (
	"errors"
	"github.com/go-park-mail-ru/2020_1_k-on/application/film"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	person "github.com/go-park-mail-ru/2020_1_k-on/application/person"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"strconv"
)

type FilmHandler struct {
	fusecase film.Usecase
	pusecase person.UseCase
}

func NewFilmHandler(e *echo.Echo, fusecase film.Usecase, pusecase person.UseCase) {
	handler := &FilmHandler{
		fusecase: fusecase,
		pusecase: pusecase,
	}
	e.GET("/films/:id", handler.GetFilm)
	e.GET("/", handler.GetFilmList)
	e.GET("/films", handler.FilterFilmList)
	e.GET("/films/filter", handler.FilterFilmData)
	e.POST("/films", handler.CreateFilm)
}

func (fh FilmHandler) FilterFilmData(ctx echo.Context) error {
	d, ok := fh.fusecase.FilterFilmData()
	if !ok {
		resp, _ := models.Generate(500, "can't get data", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("can't get data")
	}
	resp, _ := models.Generate(200, d, &ctx).MarshalJSON()
	_, err := ctx.Response().Write(resp)
	return err
}

func (fh FilmHandler) FilterFilmList(ctx echo.Context) error {
	query := ctx.QueryParams()
	f, ok := fh.fusecase.FilterFilmList(query)
	if !ok {
		resp, _ := models.Generate(500, "can't get films", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("can't get films")
	}
	var fl models.ListsFilm
	resp, _ := models.Generate(200, fl.Convert(f), &ctx).MarshalJSON()
	_, err := ctx.Response().Write(resp)
	return err
}

func (fh FilmHandler) GetFilm(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp, _ := models.Generate(400, "not number", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return err
	}
	f, ok := fh.fusecase.GetFilm(uint(id))
	if !ok {
		resp, _ := models.Generate(404, "Not Found", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("Not Found")
	}
	a, _ := fh.pusecase.GetActorsForFilm(f.ID)
	g, _ := fh.fusecase.GetFilmGenres(f.ID)
	r := make(map[string]interface{})
	r["object"] = f
	r["actors"] = a
	r["genres"] = g
	resp, _ := models.Generate(200, r, &ctx).MarshalJSON()
	_, err = ctx.Response().Write(resp)
	return err
}

func (fh FilmHandler) CreateFilm(ctx echo.Context) error {
	film := models.Film{}
	defer ctx.Request().Body.Close()
	err := easyjson.UnmarshalFromReader(ctx.Request().Body, &film)
	if err != nil || film.EnglishName == "" {
		resp, _ := models.Generate(400, "request parser error", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return err
	}
	f, ok := fh.fusecase.CreateFilm(film)
	if !ok {
		resp, _ := models.Generate(500, "can't create film", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("can't create film")
	}
	resp, _ := models.Generate(200, f, &ctx).MarshalJSON()
	ctx.Response().Write(resp)
	return err
}

func (fh FilmHandler) GetFilmList(ctx echo.Context) error {
	f, ok := fh.fusecase.GetFilmsList(13, 0)
	if !ok {
		resp, _ := models.Generate(500, "can't get films", &ctx).MarshalJSON()
		ctx.Response().Write(resp)
		return errors.New("can't get films")
	}
	var fl models.ListsFilm
	r := make(map[string]interface{})
	rec := fl.Convert(f[0:5])
	for ind, _ := range rec {
		rec[ind].Image = f[ind].BackgroundImage
	}
	r["recommendations"] = rec
	coll := make([]models.Collection, 2)
	coll[0] = models.Collection{"Сейчас смотрят", fl.Convert(f)}
	coll[1] = models.Collection{"Новое", fl.Convert(f)}
	r["collections"] = coll
	resp, _ := models.Generate(200, r, &ctx).MarshalJSON()
	_, err := ctx.Response().Write(resp)
	return err
}
