package usecase

import (
	mockseries "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/filter/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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
