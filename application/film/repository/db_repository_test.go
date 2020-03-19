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

func TestPostgresForFilms_GetById(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	var elemID uint = 1
	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name", "agelimit", "image"})
	expect := models.Film{elemID, "name", 10, "image"}
	rows = rows.AddRow(expect.ID, expect.Name, expect.AgeLimit, expect.Image)
	mock.ExpectQuery(`SELECT id,name,agelimit,image FROM "films" WHERE.*$`).
		//WithArgs(elemID).
		WillReturnRows(rows)

	repo := &PostgresForFilms{
		DB: DB,
	}
	item, ok := repo.GetById(elemID)
	if !ok {
		t.Error(ok)
		t.Error(rows)
		t.Error(expect)
		t.Error(item)
		return
	}
	require.Equal(t, *item, expect)
	require.True(t, ok)
}

func TestPostgresForFilms_GetById2(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	var elemID uint = 1
	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name", "agelimit", "image"})
	expect := models.Film{0, "", 0, ""}
	rows = rows.AddRow(expect.ID, expect.Name, expect.AgeLimit, expect.Image)
	mock.ExpectQuery(`SELECT id,name,agelimit,image FROM "films" WHERE.*$`).
		WillReturnError(errors.New(""))

	repo := &PostgresForFilms{
		DB: DB,
	}
	item, ok := repo.GetById(elemID)
	require.Equal(t, *item, expect)
	require.False(t, ok)
}

func TestPostgresForFilms_GetByName(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	elemName := "name"
	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name", "agelimit", "image"})
	expect := models.Film{1, elemName, 10, "image"}
	rows = rows.AddRow(expect.ID, expect.Name, expect.AgeLimit, expect.Image)
	mock.ExpectQuery(`SELECT id,name,agelimit,image FROM "films" WHERE.*$`).
		//WithArgs(elemID).
		WillReturnRows(rows)

	repo := &PostgresForFilms{
		DB: DB,
	}
	item, ok := repo.GetByName(elemName)
	if !ok {
		t.Error(ok)
		t.Error(rows)
		t.Error(expect)
		t.Error(item)
		return
	}
	require.Equal(t, *item, expect)
	require.True(t, ok)
}

func TestPostgresForFilms_GetByName2(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	elemName := "name"
	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name", "agelimit", "image"})
	expect := models.Film{0, "", 0, ""}
	rows = rows.AddRow(expect.ID, expect.Name, expect.AgeLimit, expect.Image)
	mock.ExpectQuery(`SELECT id,name,agelimit,image FROM "films" WHERE.*$`).
		WillReturnError(errors.New(""))

	repo := &PostgresForFilms{
		DB: DB,
	}
	item, ok := repo.GetByName(elemName)
	require.Equal(t, *item, expect)
	require.False(t, ok)
}

func TestPostgresForFilms_GetFilmsArr(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	elemName := "name"
	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name", "agelimit", "image"})
	expect := models.Film{1, elemName, 10, "image"}
	rows = rows.AddRow(expect.ID, expect.Name, expect.AgeLimit, expect.Image)
	mock.ExpectQuery(`SELECT id,name,agelimit,image FROM "films" LIMIT.*$`).
		//WithArgs(elemID).
		WillReturnRows(rows)

	repo := &PostgresForFilms{
		DB: DB,
	}
	item, ok := repo.GetFilmsArr(0, 1)
	if !ok {
		t.Error(ok)
		t.Error(rows)
		t.Error(expect)
		t.Error(item)
		return
	}
	require.Equal(t, *item, models.Films{expect})
	require.True(t, ok)
}

func TestPostgresForFilms_GetFilmsArr2(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	elemName := "name"
	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name", "agelimit", "image"})
	expect := models.Film{1, elemName, 10, "image"}
	rows = rows.AddRow(expect.ID, expect.Name, expect.AgeLimit, expect.Image)
	mock.ExpectQuery(`SELECT id,name,agelimit,image FROM "films" LIMIT.*$`).
		WillReturnError(errors.New(""))
	repo := &PostgresForFilms{
		DB: DB,
	}
	item, ok := repo.GetFilmsArr(0, 1)
	require.Equal(t, *item, models.Films{})
	require.False(t, ok)
}

func TestPostgresForFilms_Create(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	name := "name"
	image := "image"
	agelimit := 10

	testFilm := models.Film{
		Name:     name,
		AgeLimit: agelimit,
		Image:    image,
	}

	// good query
	rows := sqlmock.
		NewRows([]string{"id"})
	rows = rows.AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO`).
		WithArgs(name, agelimit, image).WillReturnRows(rows)
	mock.ExpectCommit()

	repo := &PostgresForFilms{
		DB: DB,
	}
	item, ok := repo.Create(&testFilm)
	if !ok {
		t.Error(ok)
		t.Error(item)
		return
	}
	require.Equal(t, item.ID, uint(1))
	require.True(t, ok)
}

func TestPostgresForFilms_Create2(t *testing.T) {
	mock, DB := SetupDB()
	defer DB.Close()

	name := "name"
	image := "image"
	agelimit := 10

	testFilm := models.Film{
		Name:     name,
		AgeLimit: agelimit,
		Image:    image,
	}

	// good query
	rows := sqlmock.
		NewRows([]string{"id"})
	rows = rows.AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO`).
		WithArgs(name, agelimit, image).WillReturnRows(rows)
	mock.ExpectCommit().WillReturnError(errors.New(""))

	repo := &PostgresForFilms{
		DB: DB,
	}
	item, ok := repo.Create(&testFilm)
	require.Equal(t, item, models.Film{})
	require.False(t, ok)
}
