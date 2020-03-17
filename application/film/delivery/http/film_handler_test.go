package http

import (
	mockfilm "2020_1_k-on/application/film/mocks"
	"2020_1_k-on/application/film/usecase"
	"2020_1_k-on/application/models"
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var name = "name"
var image = "image"
var agelimit = 10
var fid = uint(1)

var testFilm = models.Film{
	ID:       fid,
	Name:     name,
	AgeLimit: agelimit,
	Image:    image,
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
	c, fh, films := setupEcho(t, "/films", http.MethodGet)
	films.EXPECT().GetFilmsArr(uint(10), uint(0)).Return(&models.Films{testFilm}, true)
	err := fh.GetFilmList(c)
	require.Equal(t, err, nil)
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
