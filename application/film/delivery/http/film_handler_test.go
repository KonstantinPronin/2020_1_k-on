package http

import (
	mockfilm "2020_1_k-on/application/film/mocks"
	"2020_1_k-on/application/film/usecase"
	"2020_1_k-on/application/models"
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

func TestFilmHandler_GetFilm(t *testing.T) {
	e := echo.New()
	r := e.Router()
	r.Add("GET", "/films/:id", func(echo.Context) error { return nil })
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/films/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := usecase.NewFilmUsecase(films)
	films.EXPECT().GetById(gomock.Eq(testFilm.ID)).Return(&testFilm, true)

	fh := FilmHandler{usecase: usecase}
	err := fh.GetFilm(c)
	require.Equal(t, err, nil)
}
