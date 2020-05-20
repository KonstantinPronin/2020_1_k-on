package usecase

import (
	mockseries "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/filter/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSerialUsecase_FilterSeriesData(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	usecase := NewSeriesFilter(series)
	series.EXPECT().FilterSeriesData().Return(nil, true)

	s, ok := usecase.FilterSeriesData()
	if !ok {
		t.Error(s)
	}
	require.True(t, ok)
}

func TestSerialUsecase_FilterSeriesData2(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	usecase := NewSeriesFilter(series)
	series.EXPECT().FilterSeriesData().Return(nil, false)

	_, ok := usecase.FilterSeriesData()
	require.False(t, ok)
}

func TestSerialUsecase_FilterSeriesList(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	usecase := NewSeriesFilter(series)
	series.EXPECT().FilterSeriesList(nil).Return(&models.SeriesArr{}, true)

	s, ok := usecase.FilterSeriesList(nil)
	if !ok {
		t.Error(s)
	}
	require.True(t, ok)
}

func TestSerialUsecase_FilterSeriesList2(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	usecase := NewSeriesFilter(series)
	series.EXPECT().FilterSeriesList(nil).Return(&models.SeriesArr{}, false)

	_, ok := usecase.FilterSeriesList(nil)
	require.False(t, ok)
}
