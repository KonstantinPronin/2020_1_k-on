package usecase

import (
	mockfilm "github.com/go-park-mail-ru/2020_1_k-on/application/film/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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

func TestFilmUsecase_GetFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := NewFilmUsecase(films)
	films.EXPECT().GetById(gomock.Eq(testFilm.ID)).Return(&testFilm, true)

	f, ok := usecase.GetFilm(testFilm.ID)
	if !ok {
		t.Error(f)
	}
	require.Equal(t, testFilm, f)
	require.True(t, ok)
}

func TestFilmUsecase_GetFilm2(t *testing.T) {
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := NewFilmUsecase(films)
	films.EXPECT().GetById(gomock.Eq(testFilm.ID)).Return(&models.Film{}, false)

	f, ok := usecase.GetFilm(testFilm.ID)
	require.Equal(t, models.Film{}, f)
	require.False(t, ok)

}

func TestFilmUsecase_GetFilmsList(t *testing.T) {
	tfilms := models.Films{testFilm}
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := NewFilmUsecase(films)
	films.EXPECT().GetFilmsArr(uint(10), uint(0)).Return(&tfilms, true)

	f, ok := usecase.GetFilmsList()
	require.Equal(t, tfilms, f)
	require.True(t, ok)
}

func TestFilmUsecase_GetFilmsList2(t *testing.T) {
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := NewFilmUsecase(films)
	films.EXPECT().GetFilmsArr(uint(10), uint(0)).Return(&models.Films{}, false)
	f, ok := usecase.GetFilmsList()
	require.Equal(t, models.Films{}, f)
	require.False(t, ok)
}

func TestFilmUsecase_CreateFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := NewFilmUsecase(films)
	films.EXPECT().Create(&testFilm).Return(testFilm, true)

	f, ok := usecase.CreateFilm(testFilm)
	if !ok {
		t.Error(f)
	}
	require.Equal(t, testFilm, f)
	require.True(t, ok)
}
