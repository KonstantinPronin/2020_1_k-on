package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

var testPerson = models.Person{
	Id:         1,
	Name:       "test",
	Occupation: "test",
	BirthDate:  "2020-07-22",
	BirthPlace: "test",
	Image:      "test",
	Films:      nil,
	Series:     nil,
}

var testListPerson = models.ListPerson{
	Id:   1,
	Name: "test",
}

var testFilm = models.ListFilm{
	ID:          1,
	MainGenre:   "test",
	RussianName: "test",
	Image:       "test",
	Country:     "test",
	Year:        0,
	AgeLimit:    0,
	Rating:      0,
}

var testSeries = models.ListSeries{
	ID:          0,
	MainGenre:   "test",
	RussianName: "test",
	Image:       "test",
	Country:     "test",
	YearFirst:   0,
	YearLast:    0,
	AgeLimit:    0,
	Rating:      0,
}

func initMockDb(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	conn, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	return conn, mock
}

func TestPersonDatabase_Add(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	persons := NewPersonDatabase(db, zap.NewExample())

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*)persons`).
		WithArgs(testPerson.Id, testPerson.Name, testPerson.Occupation, testPerson.BirthDate, testPerson.BirthPlace, testPerson.Image).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	if _, err := persons.Add(&testPerson); err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPersonDatabase_Update(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	persons := NewPersonDatabase(db, zap.NewExample())

	mock.ExpectQuery(`SELECT (\*) FROM (.*)"persons" WHERE (.*)"persons"."id" (.*) LIMIT 1`).
		WithArgs(testPerson.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "occupation", "birth_date", "birth_place", "image"}).
			AddRow(testPerson.Id, "old", testPerson.Occupation, testPerson.BirthDate, testPerson.BirthPlace, testPerson.Image))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE (.*)"persons" SET (.*) WHERE (.*)"persons"`).WithArgs(
		testPerson.Name, testPerson.Occupation, testPerson.BirthDate, testPerson.BirthPlace, testPerson.Image, testPerson.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if _, err := persons.Update(&testPerson); err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPersonDatabase_GetById(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	persons := NewPersonDatabase(db, zap.NewExample())

	mock.ExpectQuery(`SELECT (\*) FROM (.*)"persons" WHERE (.*)"persons"."id" (.*) LIMIT 1`).
		WithArgs(testPerson.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "occupation", "birth_date", "birth_place", "image"}).
			AddRow(testPerson.Id, "old", testPerson.Occupation, testPerson.BirthDate, testPerson.BirthPlace, testPerson.Image))
	mock.ExpectQuery(`SELECT (.*) FROM (.*)film_actor(.*) inner join (.*)films(.*) on (.*)film_id = (.*)id ` +
		`WHERE (.*)person_id = (.*)`).WithArgs(testPerson.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "maingenre", "russianname", "image", "country", "year", "agelimit", "rating"}).
			AddRow(testFilm.ID, testFilm.MainGenre, testFilm.RussianName, testFilm.Image, testFilm.Country,
				testFilm.Year, testFilm.AgeLimit, testFilm.Rating))
	mock.ExpectQuery(`SELECT (.*) FROM (.*)series_actor(.*) inner join (.*)series(.*) on (.*)series_id = (.*)id ` +
		`WHERE (.*)person_id = (.*)`).WithArgs(testPerson.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "maingenre", "russianname", "image", "country", "yearfirst", "yearlast", "agelimit", "rating"}).
			AddRow(testSeries.ID, testSeries.MainGenre, testSeries.RussianName, testSeries.Image, testSeries.Country,
				testSeries.YearFirst, testSeries.YearLast, testSeries.AgeLimit, testSeries.Rating))

	if _, err := persons.GetById(testPerson.Id); err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPersonDatabase_GetActorsForFilm(t *testing.T) {
	db, mock := initMockDb(t)
	filmId := 1
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	persons := NewPersonDatabase(db, zap.NewExample())

	mock.ExpectQuery(`SELECT (.*) FROM (.*)film_actor(.*) inner join (.*)persons(.*) on (.*)person_id = (.*)id ` +
		`WHERE (.*)film_id = `).WithArgs(filmId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(testListPerson.Id, testListPerson.Name))

	per, err := persons.GetActorsForFilm(testPerson.Id)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	assert.Equal(t, testListPerson, per[0])

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPersonDatabase_GetActorsForSeries(t *testing.T) {
	db, mock := initMockDb(t)
	filmId := 1
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	persons := NewPersonDatabase(db, zap.NewExample())

	mock.ExpectQuery(`SELECT (.*) FROM (.*)series_actor(.*) inner join (.*)persons(.*) on (.*)person_id = (.*)id ` +
		`WHERE (.*)series_id = `).WithArgs(filmId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(testListPerson.Id, testListPerson.Name))

	per, err := persons.GetActorsForSeries(testPerson.Id)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	assert.Equal(t, testListPerson, per[0])

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
