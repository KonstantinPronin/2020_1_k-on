package http

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/client"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	person "github.com/go-park-mail-ru/2020_1_k-on/application/person"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type SeriesHandler struct {
	rpcSeriesFilter client.ISeriesFilterClient
	usecase         series.Usecase
	pusecase        person.UseCase
}

func NewSeriesHandler(e *echo.Echo,
	rpcSeriesFilter client.ISeriesFilterClient,
	usecase series.Usecase, pusecase person.UseCase) {
	handler := &SeriesHandler{
		rpcSeriesFilter: rpcSeriesFilter,
		usecase:         usecase,
		pusecase:        pusecase,
	}
	e.GET("/series/:id", handler.GetSeries)
	e.GET("/series/:id/seasons", handler.GetSeriesSeasons)
	e.GET("/seasons/:id/episodes", handler.GetSeasonEpisodes)

	e.GET("/series", handler.FilterSeriesList)
	e.GET("/series/filter", handler.FilterSeriesData)
	e.GET("/series/search/:word", handler.Search)
}

func (sh SeriesHandler) FilterSeriesData(ctx echo.Context) error {
	d, ok := sh.rpcSeriesFilter.GetFilterFields()
	if !ok {
		return middleware.WriteErrResponse(ctx,
			http.StatusInternalServerError, "can't get data")
	}
	return middleware.WriteOkResponse(ctx, d)
}

func (sh SeriesHandler) FilterSeriesList(ctx echo.Context) error {
	query := ctx.QueryParams()
	s, ok := sh.rpcSeriesFilter.GetFilteredSeries(query)
	if !ok {
		return middleware.WriteErrResponse(ctx,
			http.StatusInternalServerError, "can't get series")
	}
	var sl models.ListSeriesArr
	return middleware.WriteOkResponse(ctx, sl.Convert(s))

}

func (sh SeriesHandler) GetSeries(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return middleware.WriteErrResponse(ctx,
			http.StatusBadRequest, "not number")
	}
	serial, ok := sh.usecase.GetSeriesByID(uint(id))
	if !ok {
		return middleware.WriteErrResponse(ctx,
			http.StatusNotFound, "Not Found")

	}
	a, err := sh.pusecase.GetActorsForSeries(serial.ID)
	if err != nil {
		return middleware.WriteErrResponse(ctx,
			http.StatusBadRequest, "not number")
	}
	g, _ := sh.usecase.GetSeriesGenres(serial.ID)
	simFilms, _ := sh.usecase.GetSimilarFilms(uint(id))
	var fl models.ListsFilm
	simSeries, _ := sh.usecase.GetSimilarSeries(uint(id))
	var sl models.ListSeriesArr

	r := make(map[string]interface{})
	r["object"] = serial
	r["actors"] = a
	r["genres"] = g
	r["simfilms"] = fl.Convert(simFilms)
	r["simseries"] = sl.Convert(simSeries)
	return middleware.WriteOkResponse(ctx, r)

}

func (sh SeriesHandler) GetSeriesSeasons(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return middleware.WriteErrResponse(ctx,
			http.StatusBadRequest, "not number")
	}
	seasons, ok := sh.usecase.GetSeriesSeasons(uint(id))
	if !ok {
		return middleware.WriteErrResponse(ctx,
			http.StatusNotFound, "Not Found")
	}
	return middleware.WriteOkResponse(ctx, seasons)

}

func (sh SeriesHandler) GetSeasonEpisodes(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return middleware.WriteErrResponse(ctx,
			http.StatusBadRequest, "not number")
	}
	episodes, ok := sh.usecase.GetSeasonEpisodes(uint(id))
	if !ok {
		return middleware.WriteErrResponse(ctx,
			http.StatusNotFound, "Not Found")
	}

	return middleware.WriteOkResponse(ctx, episodes)
}

func (sh SeriesHandler) Search(ctx echo.Context) error {
	word := ctx.Param("word")
	query := ctx.QueryParams()

	series, ok := sh.usecase.Search(word, query)
	if !ok {
		return middleware.WriteErrResponse(ctx, http.StatusInternalServerError, "can't find series")
	}

	var sl models.ListSeriesArr
	return middleware.WriteOkResponse(ctx, sl.Convert(series))
}
