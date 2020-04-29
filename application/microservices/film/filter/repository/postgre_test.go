package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

var name = "name"
var reference = "reference"
var testGenre = models.Genre{
	Name:      name,
	Reference: reference,
}

var image = "image"
var ftype1 = "film"

//var ftype2 = "series"
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
var year = 2012
var agelimit = 10
var fid = uint(1)

var testFilm = models.Film{
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
	Year:            year,
	AgeLimit:        agelimit,
	SumVotes:        sumvotes,
	TotalVotes:      totalvotes,
	BackgroundImage: image,
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

func TestPostgresForFilms_FilterFilmsList(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	// good query

	rows := sqlmock.
		NewRows([]string{"id", "maingenre", "russianname", "englishname", "trailerlink",
			"rating", "imdbrating", "totalvotes", "sumvotes", "description", "image", "backgroundimage",
			"country", "year", "agelimit"})
	expect := models.Film(testFilm)
	rows = rows.AddRow(expect.ID, expect.MainGenre, expect.RussianName, expect.EnglishName,
		expect.TrailerLink, expect.Rating, expect.ImdbRating, expect.TotalVotes, expect.SumVotes,
		expect.Description, expect.Image, expect.BackgroundImage, expect.Country, expect.Year, expect.AgeLimit)
	rows = rows.AddRow(expect.ID+1, expect.MainGenre, expect.RussianName, expect.EnglishName,
		expect.TrailerLink, expect.Rating, expect.ImdbRating, expect.TotalVotes, expect.SumVotes,
		expect.Description, expect.Image, expect.BackgroundImage, expect.Country, expect.Year, expect.AgeLimit)
	mock.ExpectQuery(`SELECT (.*)(\*) FROM (.*)"films" (.*)`).
		WillReturnRows(rows)

	repo := &FilmFiltersDb{
		DB: DB,
	}
	query := make(map[string][]string)
	query["year"] = []string{"2012"}
	query["page"] = []string{"1"}
	item, ok := repo.FilterFilmsList(query)
	if !ok {
		t.Error(ok)
		t.Error(rows)
		t.Error(expect)
		t.Error(item)
		return
	}
	expect2 := expect
	expect2.ID += 1
	require.Equal(t, *item, models.Films{expect, expect2})
	require.True(t, ok)
}

func TestPostgresForFilms_FilterFilmsList2(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	// good query

	rows := sqlmock.
		NewRows([]string{"id", "maingenre", "russianname", "englishname", "trailerlink",
			"rating", "imdbrating", "totalvotes", "sumvotes", "description", "image", "backgroundimage",
			"country", "year", "agelimit"})
	expect := models.Film(testFilm)
	rows = rows.AddRow(expect.ID, expect.MainGenre, expect.RussianName, expect.EnglishName,
		expect.TrailerLink, expect.Rating, expect.ImdbRating, expect.TotalVotes, expect.SumVotes,
		expect.Description, expect.Image, expect.BackgroundImage, expect.Country, expect.Year, expect.AgeLimit)
	rows = rows.AddRow(expect.ID+1, expect.MainGenre, expect.RussianName, expect.EnglishName,
		expect.TrailerLink, expect.Rating, expect.ImdbRating, expect.TotalVotes, expect.SumVotes,
		expect.Description, expect.Image, expect.BackgroundImage, expect.Country, expect.Year, expect.AgeLimit)
	mock.ExpectQuery(`SELECT (.*)(\*) FROM (.*)"films" (.*)`).
		WillReturnError(errors.New(""))

	repo := &FilmFiltersDb{
		DB: DB,
	}
	query := make(map[string][]string)
	query["year"] = []string{"ALL"}
	query["order"] = []string{"rating"}
	query["genre"] = []string{"1"}
	query["page"] = []string{"1"}
	item, ok := repo.FilterFilmsList(query)
	expect2 := expect
	expect2.ID += 1
	require.NotEqual(t, *item, models.Films{expect, expect2})
	require.False(t, ok)
}

func TestPostgresForFilms_FilterFilmData(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	// good query
	rows2 := sqlmock.
		NewRows([]string{"name", "reference"})
	expect2 := testGenre
	rows2 = rows2.AddRow(expect2.Name, expect2.Reference)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"genres" `).
		WillReturnRows(rows2)

	rows := sqlmock.
		NewRows([]string{"max", "min"})
	expect := models.Film(testFilm)
	rows = rows.AddRow(expect.Year+1, expect.Year)
	mock.ExpectQuery(`SELECT (.*)" `).
		WillReturnRows(rows)

	repo := NewFilmFiltersDb(DB)
	item, ok := repo.FilterFilmData()
	if !ok {
		t.Error(ok)
		t.Error(rows)
		t.Error(expect)
		t.Error(item)
		return
	}
	resp := make(map[string]models.Genres)
	resp["genre"] = models.Genres{models.Genre{"Все жанры", "%"}, expect2}
	resp["year"] = models.Genres{
		models.Genre{"Все годы", "%"},
	}

	require.Equal(t, item["genres"], resp["genres"])
	require.Equal(t, item["filters"], resp["filters"])
	require.True(t, ok)
}

func TestPostgresForFilms_FilterFilmData2(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	// good query
	rows2 := sqlmock.
		NewRows([]string{"name", "reference"})
	expect2 := testGenre
	rows2 = rows2.AddRow(expect2.Name, expect2.Reference)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"genres" `).
		WillReturnError(errors.New(""))
	rows := sqlmock.
		NewRows([]string{"max", "min"})
	expect := models.Film(testFilm)
	rows = rows.AddRow(expect.Year+1, expect.Year)
	mock.ExpectQuery(`SELECT (.*)" `).
		WillReturnRows(rows)

	repo := NewFilmFiltersDb(DB)
	_, ok := repo.FilterFilmData()
	require.False(t, ok)
}

func TestPostgresForFilms_FilterFilmData3(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	// good query
	rows2 := sqlmock.
		NewRows([]string{"name", "reference"})
	expect2 := testGenre
	rows2 = rows2.AddRow(expect2.Name, expect2.Reference)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"genres" `).
		WillReturnRows(rows2)

	rows := sqlmock.
		NewRows([]string{"max", "min"})
	expect := models.Film(testFilm)
	rows = rows.AddRow(expect.Year+1, expect.Year)
	mock.ExpectQuery(`SELECT (.*)" `).
		WillReturnError(errors.New(""))

	repo := NewFilmFiltersDb(DB)
	_, ok := repo.FilterFilmData()

	require.False(t, ok)
}
