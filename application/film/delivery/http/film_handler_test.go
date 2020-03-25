package http

import (
	"bytes"
	mockfilm "github.com/go-park-mail-ru/2020_1_k-on/application/film/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/film/usecase"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var image = "image"
var ftype1 = "film"

//var ftype2 = "serial"
var mg = "mg"
var rn = "rn"
var en = "en"
var seasons = 1
var tl = "tl"
var rating = 1.2
var imdbrating = 9.87
var d = "d"
var c = "c"
var year = 2012
var agelimit = 10
var fid = uint(1)

var testFilm = models.Film{
	ID:          fid,
	Type:        ftype1,
	MainGenre:   mg,
	RussianName: rn,
	EnglishName: en,
	Seasons:     seasons,
	TrailerLink: tl,
	Rating:      rating,
	ImdbRating:  imdbrating,
	Description: d,
	Image:       image,
	Country:     c,
	Year:        year,
	AgeLimit:    agelimit,
}

func setupEcho(t *testing.T, url, method string) (echo.Context, FilmHandler, *mockfilm.MockRepository) {
	e := echo.New()
	r := e.Router()
	r.Add(method, url, func(echo.Context) error { return nil })
	var req *http.Request
	switch method {
	case http.MethodPost:
		f, _ := testFilm.MarshalJSON()
		req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(f))
	case http.MethodGet:
		req = httptest.NewRequest(http.MethodGet, url, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(url)
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := usecase.NewFilmUsecase(films)
	fh := FilmHandler{usecase: usecase}
	return c, fh, films

}

func TestFilmHandler_GetFilm(t *testing.T) {
	c, fh, films := setupEcho(t, "/films/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	films.EXPECT().GetById(gomock.Eq(testFilm.ID)).Return(&testFilm, true)
	err := fh.GetFilm(c)
	require.Equal(t, err, nil)
}

func TestFilmHandler_GetFilmList(t *testing.T) {
	c, fh, films := setupEcho(t, "/", http.MethodGet)
	films.EXPECT().GetFilmsArr(uint(10), uint(0)).Return(&models.Films{testFilm}, true)
	err := fh.GetFilmList(c)
	require.Equal(t, err, nil)
}

func TestFilmHandler_GetFilmList2(t *testing.T) {
	c, fh, films := setupEcho(t, "/", http.MethodGet)
	films.EXPECT().GetFilmsArr(uint(10), uint(0)).Return(&models.Films{}, false)
	err := fh.GetFilmList(c)
	require.NotEqual(t, err, nil)
}

func TestFilmHandler_CreateFilm(t *testing.T) {
	c, fh, films := setupEcho(t, "/films", http.MethodPost)
	films.EXPECT().Create(&testFilm).Return(testFilm, true)
	err := fh.CreateFilm(c)
	require.Equal(t, err, nil)
}

func TestFilmHandler_CreateFilm2(t *testing.T) {
	c, fh, films := setupEcho(t, "/films", http.MethodGet)
	films.EXPECT().Create(&testFilm).Return(testFilm, true)
	err := fh.CreateFilm(c)
	require.NotEqual(t, err, nil)
}

func TestFilmHandler_CreateFilm3(t *testing.T) {
	c, fh, films := setupEcho(t, "/films", http.MethodPost)
	films.EXPECT().Create(&testFilm).Return(models.Film{}, false)
	err := fh.CreateFilm(c)
	require.NotEqual(t, err, nil)
}

func TestFilmHandler_GetFilm2(t *testing.T) {
	c, fh, films := setupEcho(t, "/films/:id", http.MethodGet)
	films.EXPECT().GetById(gomock.Eq(testFilm.ID)).Return(&testFilm, true)
	err := fh.GetFilm(c)
	require.NotEqual(t, err, nil)
}

func TestFilmHandler_GetFilm3(t *testing.T) {
	c, fh, films := setupEcho(t, "/films/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	films.EXPECT().GetById(testFilm.ID).Return(&models.Film{}, false)
	err := fh.GetFilm(c)
	require.NotEqual(t, err, nil)
}

func TestFilmHandler_FilterFilmData(t *testing.T) {
	c, fh, films := setupEcho(t, "/films/filter", http.MethodGet)
	films.EXPECT().FilterFilmData().Return(nil, true)
	err := fh.FilterFilmData(c)
	require.Equal(t, err, nil)
}

func TestFilmHandler_FilterFilmData2(t *testing.T) {
	c, fh, films := setupEcho(t, "/films/filter", http.MethodGet)
	films.EXPECT().FilterFilmData().Return(nil, false)
	err := fh.FilterFilmData(c)
	require.NotEqual(t, err, nil)
}

func TestFilmHandler_FilterFilmList(t *testing.T) {
	q := make(map[string][]string)
	q["year"] = []string{"year"}
	c, fh, films := setupEcho(t, "/films", http.MethodGet)
	c.QueryParams().Add("year", "year")
	films.EXPECT().FilterFilmsList(q).Return(&models.Films{}, true)
	err := fh.FilterFilmList(c)
	require.Equal(t, err, nil)
}

func TestFilmHandler_FilterFilmList2(t *testing.T) {
	q := make(map[string][]string)
	q["year"] = []string{"year"}
	c, fh, films := setupEcho(t, "/films", http.MethodGet)
	c.QueryParams().Add("year", "year")
	films.EXPECT().FilterFilmsList(q).Return(&models.Films{}, false)
	err := fh.FilterFilmList(c)
	require.NotEqual(t, err, nil)
}
