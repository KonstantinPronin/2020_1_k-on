package http

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	mock_p "github.com/go-park-mail-ru/2020_1_k-on/application/person/mocks"
	usecase2 "github.com/go-park-mail-ru/2020_1_k-on/application/person/usecase"
	mockseries "github.com/go-park-mail-ru/2020_1_k-on/application/series/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series/usecase"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var image = "image"

var mg = "mg"
var rn = "rn"
var en = "en"
var sumvotes = 0
var totalvotes = 0
var tl = "tl"
var rating = 1.2
var imdbrating = 9.87
var d = "d"
var c = "c"
var yearfirst = 2012
var yearlast = yearfirst + 1
var agelimit = 10
var fid = uint(1)
var number = 1

var testSeries = models.Series{
	ID:              fid,
	MainGenre:       mg,
	RussianName:     rn,
	EnglishName:     en,
	TrailerLink:     tl,
	Rating:          rating,
	ImdbRating:      imdbrating,
	Description:     d,
	Image:           image,
	Country:         c,
	YearFirst:       yearfirst,
	YearLast:        yearlast,
	AgeLimit:        agelimit,
	SumVotes:        sumvotes,
	TotalVotes:      totalvotes,
	BackgroundImage: image,
}

var testSeason = models.Season{
	ID:          fid,
	SeriesID:    fid,
	Name:        rn,
	Number:      number,
	TrailerLink: tl,
	Description: d,
	Year:        yearfirst,
	Image:       image,
}

var testEpisode = models.Episode{
	ID:       fid,
	SeasonId: fid,
	Name:     rn,
	Number:   number,
	Image:    image,
}

func setupEcho(t *testing.T, url, method string) (echo.Context, SeriesHandler, *mockseries.MockRepository, *mock_p.MockRepository) {
	e := echo.New()
	r := e.Router()
	r.Add(method, url, func(echo.Context) error { return nil })
	var req *http.Request
	switch method {
	case http.MethodGet:
		req = httptest.NewRequest(http.MethodGet, url, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(url)
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	seriesuc := usecase.NewSeriesUsecase(series)
	person := mock_p.NewMockRepository(ctrl)
	persons := usecase2.NewPerson(person, nil)
	sh := SeriesHandler{usecase: seriesuc, pusecase: persons}
	return c, sh, series, person

}

func TestSeriesHandler_GetSeries(t *testing.T) {
	c, sh, series, person := setupEcho(t, "/series/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	series.EXPECT().GetSeriesByID(gomock.Eq(testSeries.ID)).Return(testSeries, true)
	series.EXPECT().GetSeriesGenres(testSeries.ID).Return(nil, true)
	person.EXPECT().GetActorsForSeries(testSeries.ID).Return(nil, nil)
	err := sh.GetSeries(c)
	require.Equal(t, err, nil)
}

func TestSeriesHandler_GetSeries2(t *testing.T) {
	c, sh, series, person := setupEcho(t, "/series/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("abc")
	series.EXPECT().GetSeriesByID(gomock.Eq(testSeries.ID)).Return(testSeries, true)
	series.EXPECT().GetSeriesGenres(testSeries.ID).Return(nil, true)
	person.EXPECT().GetActorsForSeries(testSeries.ID).Return(nil, nil)
	err := sh.GetSeries(c)
	require.NotEqual(t, err, nil)
}

func TestSeriesHandler_GetSeries3(t *testing.T) {
	c, sh, series, person := setupEcho(t, "/series/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	series.EXPECT().GetSeriesByID(gomock.Eq(testSeries.ID)).Return(testSeries, false)
	series.EXPECT().GetSeriesGenres(testSeries.ID).Return(nil, true)
	person.EXPECT().GetActorsForSeries(testSeries.ID).Return(nil, nil)
	err := sh.GetSeries(c)
	require.NotEqual(t, err, nil)
}

func TestSeriesHandler_GetSeriesSeasons(t *testing.T) {
	c, sh, series, _ := setupEcho(t, "/series/:id/seasons", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	series.EXPECT().GetSeriesSeasons(gomock.Eq(testSeries.ID)).Return(models.Seasons{testSeason}, true)
	err := sh.GetSeriesSeasons(c)
	require.Equal(t, err, nil)
}

func TestSeriesHandler_GetSeriesSeasons2(t *testing.T) {
	c, sh, series, _ := setupEcho(t, "/series/:id/seasons", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("abc")
	series.EXPECT().GetSeriesSeasons(gomock.Eq(testSeries.ID)).Return(models.Seasons{testSeason}, true)
	err := sh.GetSeriesSeasons(c)
	require.NotEqual(t, err, nil)
}

func TestSeriesHandler_GetSeriesSeasons3(t *testing.T) {
	c, sh, series, _ := setupEcho(t, "/series/:id/seasons", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	series.EXPECT().GetSeriesSeasons(gomock.Eq(testSeries.ID)).Return(models.Seasons{testSeason}, false)
	err := sh.GetSeriesSeasons(c)
	require.NotEqual(t, err, nil)
}

func TestSeriesHandler_GetSeasonEpisodes(t *testing.T) {
	c, sh, series, _ := setupEcho(t, "/seasons/:id/series", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	series.EXPECT().GetSeasonEpisodes(gomock.Eq(testSeason.ID)).Return(models.Episodes{testEpisode}, true)
	err := sh.GetSeasonEpisodes(c)
	require.Equal(t, err, nil)
}

func TestSeriesHandler_GetSeasonEpisodes2(t *testing.T) {
	c, sh, series, _ := setupEcho(t, "/seasons/:id/series", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("abc")
	series.EXPECT().GetSeasonEpisodes(gomock.Eq(testSeason.ID)).Return(models.Episodes{testEpisode}, true)
	err := sh.GetSeasonEpisodes(c)
	require.NotEqual(t, err, nil)
}

func TestSeriesHandler_GetSeasonEpisodes3(t *testing.T) {
	c, sh, series, _ := setupEcho(t, "/seasons/:id/series", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	series.EXPECT().GetSeasonEpisodes(gomock.Eq(testSeason.ID)).Return(models.Episodes{testEpisode}, false)
	err := sh.GetSeasonEpisodes(c)
	require.NotEqual(t, err, nil)
}

func TestSeriesHandler_FilterSeriesData(t *testing.T) {
	c, sh, series, _ := setupEcho(t, "/series/filter", http.MethodGet)
	series.EXPECT().FilterSeriesData().Return(nil, true)
	err := sh.FilterSeriesData(c)
	require.Equal(t, err, nil)
}

func TestSeriesHandler_FilterSeriesData2(t *testing.T) {
	c, sh, series, _ := setupEcho(t, "/series/filter", http.MethodGet)
	series.EXPECT().FilterSeriesData().Return(nil, false)
	err := sh.FilterSeriesData(c)
	require.NotEqual(t, err, nil)
}

func TestSeriesHandler_FilterSeriesList(t *testing.T) {
	q := make(map[string][]string)
	q["year"] = []string{"year"}
	c, sh, series, _ := setupEcho(t, "/series", http.MethodGet)
	c.QueryParams().Add("year", "year")
	series.EXPECT().FilterSeriesList(q).Return(&models.SeriesArr{}, true)
	err := sh.FilterSeriesList(c)
	require.Equal(t, err, nil)
}

func TestSeriesHandler_FilterSeriesList2(t *testing.T) {
	q := make(map[string][]string)
	q["year"] = []string{"year"}
	c, sh, series, _ := setupEcho(t, "/series", http.MethodGet)
	c.QueryParams().Add("year", "year")
	series.EXPECT().FilterSeriesList(q).Return(&models.SeriesArr{}, false)
	err := sh.FilterSeriesList(c)
	require.NotEqual(t, err, nil)
}
