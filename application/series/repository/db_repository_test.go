package repository

import (
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"log"
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

func SetupDB() (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("cant create mock: %s", err)
	}
	DB, erro := gorm.Open("postgres", db)
	if erro != nil {
		log.Fatalf("Got an unexpected error: %s", err)

	}
	return mock, DB
}

func TestPostgresForSerials_GetSeriesByID(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	rows := sqlmock.
		NewRows([]string{"id", "maingenre", "russianname", "englishname", "trailerlink",
			"rating", "imdbrating", "totalvotes", "sumvotes", "description", "image", "backgroundimage",
			"country", "yearfirst", "yearlast", "agelimit"})
	expect := testSeries
	rows = rows.AddRow(expect.ID, expect.MainGenre, expect.RussianName, expect.EnglishName,
		expect.TrailerLink, expect.Rating, expect.ImdbRating, expect.TotalVotes, expect.SumVotes,
		expect.Description, expect.Image, expect.BackgroundImage, expect.Country, expect.YearFirst, expect.YearLast, expect.AgeLimit)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"series" WHERE (.*)"series"."id" (.*)`).
		WillReturnRows(rows)

	repo := PostgresForSerials{DB: DB}
	item, ok := repo.GetSeriesByID(fid)
	if !ok {
		t.Error(ok)
		t.Error(rows)
		t.Error(expect)
		t.Error(item)
		return
	}
	require.Equal(t, item, expect)
	require.True(t, ok)
}

func TestPostgresForSerials_GetSeriesByID2(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	rows := sqlmock.
		NewRows([]string{"id", "maingenre", "russianname", "englishname", "trailerlink",
			"rating", "imdbrating", "totalvotes", "sumvotes", "description", "image", "backgroundimage",
			"country", "yearfirst", "yearlast", "agelimit"})
	expect := testSeries
	rows = rows.AddRow(expect.ID, expect.MainGenre, expect.RussianName, expect.EnglishName,
		expect.TrailerLink, expect.Rating, expect.ImdbRating, expect.TotalVotes, expect.SumVotes,
		expect.Description, expect.Image, expect.BackgroundImage, expect.Country, expect.YearFirst, expect.YearLast, expect.AgeLimit)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"series" WHERE (.*)"series"."id" (.*)`).
		WillReturnError(errors.New(""))

	repo := PostgresForSerials{DB: DB}
	item, ok := repo.GetSeriesByID(fid)
	require.NotEqual(t, item, expect)
	require.False(t, ok)
}

func TestPostgresForSerials_GetSeriesSeasons(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	rows := sqlmock.
		NewRows([]string{"id", "seriesid", "name", "trailerlink",
			"number", "description", "image", "year"})
	expect := testSeason
	rows = rows.AddRow(expect.ID, expect.SeriesID, expect.Name,
		expect.TrailerLink, expect.Number, expect.Description, expect.Image, expect.Year)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"seasons" WHERE (.*)`).
		WillReturnRows(rows)

	repo := PostgresForSerials{DB: DB}
	item, ok := repo.GetSeriesSeasons(expect.SeriesID)
	if !ok {
		t.Error(ok)
		t.Error(rows)
		t.Error(expect)
		t.Error(item)
		return
	}
	require.Equal(t, item, models.Seasons{expect})
	require.True(t, ok)
}

func TestPostgresForSerials_GetSeriesSeasons2(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	rows := sqlmock.
		NewRows([]string{"id", "seriesid", "name", "trailerlink",
			"number", "description", "image", "year"})
	expect := testSeason
	rows = rows.AddRow(expect.ID, expect.SeriesID, expect.Name,
		expect.TrailerLink, expect.Number, expect.Description, expect.Image, expect.Year)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"seasons" WHERE (.*)`).
		WillReturnError(errors.New(""))

	repo := PostgresForSerials{DB: DB}
	item, ok := repo.GetSeriesSeasons(expect.SeriesID)
	require.NotEqual(t, item, models.Seasons{expect})
	require.False(t, ok)
}

func TestPostgresForSerials_GetSeasonEpisodes(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	rows := sqlmock.
		NewRows([]string{"id", "seasonid", "name", "number", "image"})
	expect := testEpisode
	rows = rows.AddRow(expect.ID, expect.SeasonId, expect.Name,
		expect.Number, expect.Image)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"episodes" WHERE (.*)`).
		WillReturnRows(rows)

	repo := PostgresForSerials{DB: DB}
	item, ok := repo.GetSeasonEpisodes(testEpisode.SeasonId)
	if !ok {
		t.Error(ok)
		t.Error(rows)
		t.Error(expect)
		t.Error(item)
		return
	}
	require.Equal(t, item, models.Episodes{expect})
	require.True(t, ok)
}

func TestPostgresForSerials_GetSeasonEpisodes2(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	rows := sqlmock.
		NewRows([]string{"id", "seasonid", "name", "number", "image"})
	expect := testEpisode
	rows = rows.AddRow(expect.ID, expect.SeasonId, expect.Name,
		expect.Number, expect.Image)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"episodes" WHERE (.*)`).
		WillReturnError(errors.New(""))

	repo := PostgresForSerials{DB: DB}
	item, ok := repo.GetSeasonEpisodes(testEpisode.SeasonId)
	require.NotEqual(t, item, models.Episodes{expect})
	require.False(t, ok)
}
