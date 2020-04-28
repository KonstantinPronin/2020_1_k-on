package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	mockseries "github.com/go-park-mail-ru/2020_1_k-on/application/series/mocks"
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

func TestSerialUsecase_GetSeriesByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	seriesuc := NewSeriesUsecase(series)
	series.EXPECT().GetSeriesByID(testSeries.ID).Return(testSeries, true)

	ser, ok := seriesuc.GetSeriesByID(testSeries.ID)
	if !ok {
		t.Error(ser)
	}
	require.Equal(t, testSeries, ser)
	require.True(t, ok)
}

func TestSerialUsecase_GetSeriesByID2(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	seriesuc := NewSeriesUsecase(series)
	series.EXPECT().GetSeriesByID(testSeries.ID).Return(models.Series{}, false)

	ser, ok := seriesuc.GetSeriesByID(testSeries.ID)
	require.NotEqual(t, testSeries, ser)
	require.False(t, ok)
}

func TestSerialUsecase_GetSeriesSeasons(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	seriesuc := NewSeriesUsecase(series)
	series.EXPECT().GetSeriesSeasons(testSeries.ID).Return(models.Seasons{testSeason}, true)

	sesons, ok := seriesuc.GetSeriesSeasons(testSeries.ID)
	if !ok {
		t.Error(sesons)
	}
	require.Equal(t, models.Seasons{testSeason}, sesons)
	require.True(t, ok)
}

func TestSerialUsecase_GetSeriesSeasons2(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	seriesuc := NewSeriesUsecase(series)
	series.EXPECT().GetSeriesSeasons(testSeries.ID).Return(models.Seasons{}, false)

	sesons, ok := seriesuc.GetSeriesSeasons(testSeries.ID)
	require.Equal(t, models.Seasons{}, sesons)
	require.False(t, ok)
}

func TestSerialUsecase_GetSeasonEpisodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	seriesuc := NewSeriesUsecase(series)
	series.EXPECT().GetSeasonEpisodes(testSeason.ID).Return(models.Episodes{testEpisode}, true)

	episodes, ok := seriesuc.GetSeasonEpisodes(testSeason.ID)
	if !ok {
		t.Error(episodes)
	}
	require.Equal(t, models.Episodes{testEpisode}, episodes)
	require.True(t, ok)
}

func TestSerialUsecase_GetSeasonEpisodes2(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	seriesuc := NewSeriesUsecase(series)
	series.EXPECT().GetSeasonEpisodes(testSeason.ID).Return(models.Episodes{}, false)

	episodes, ok := seriesuc.GetSeasonEpisodes(testSeason.ID)
	require.Equal(t, models.Episodes{}, episodes)
	require.False(t, ok)
}

func TestSerialUsecase_GetSeriesGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	usecase := NewSeriesUsecase(series)
	series.EXPECT().GetSeriesGenres(fid).Return(nil, true)

	s, ok := usecase.GetSeriesGenres(fid)
	if !ok {
		t.Error(s)
	}
	require.True(t, ok)
}

func TestSerialUsecase_GetSeriesGenres2(t *testing.T) {
	ctrl := gomock.NewController(t)
	series := mockseries.NewMockRepository(ctrl)
	usecase := NewSeriesUsecase(series)
	series.EXPECT().GetSeriesGenres(fid).Return(nil, false)

	_, ok := usecase.GetSeriesGenres(fid)
	require.False(t, ok)
}
