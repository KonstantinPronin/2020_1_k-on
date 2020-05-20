package usecase

import (
	mockfilm "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/filter/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilmUsecase_FilterFilmList(t *testing.T) {
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := NewFilmFilter(films)
	films.EXPECT().FilterFilmsList(nil).Return(&models.Films{}, true)

	f, ok := usecase.FilterFilmList(nil)
	if !ok {
		t.Error(f)
	}
	require.True(t, ok)
}

func TestFilmUsecase_FilterFilmList2(t *testing.T) {
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := NewFilmFilter(films)
	films.EXPECT().FilterFilmsList(nil).Return(&models.Films{}, false)

	_, ok := usecase.FilterFilmList(nil)
	require.False(t, ok)
}

func TestFilmUsecase_FilterFilmData(t *testing.T) {
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := NewFilmFilter(films)
	films.EXPECT().FilterFilmData().Return(nil, true)

	f, ok := usecase.FilterFilmData()
	if !ok {
		t.Error(f)
	}
	require.True(t, ok)
}

func TestFilmUsecase_FilterFilmData2(t *testing.T) {
	ctrl := gomock.NewController(t)
	films := mockfilm.NewMockRepository(ctrl)
	usecase := NewFilmFilter(films)
	films.EXPECT().FilterFilmData().Return(nil, false)

	_, ok := usecase.FilterFilmData()
	require.False(t, ok)
}
