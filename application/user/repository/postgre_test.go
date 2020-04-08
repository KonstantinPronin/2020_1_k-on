package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

var testUser = models.User{
	Id:       1,
	Username: "test",
	Password: "test",
	Email:    "test@example.com",
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

func TestUserDatabase_Add(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	ud := NewUserDatabase(db, zap.NewExample())

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*)"users"`).WithArgs(
		testUser.Id, testUser.Username, testUser.Password, testUser.Email, testUser.Image).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(testUser.Id))
	mock.ExpectCommit()

	if err := ud.Add(&testUser); err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserDatabase_Update(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	ud := NewUserDatabase(db, zap.NewExample())

	mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
		testUser.Id).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "username", "password", "email", "image"}).AddRow(
		testUser.Id, "old", testUser.Password, testUser.Email, testUser.Image))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE (.*)"users" SET (.*) WHERE (.*)"users"`).WithArgs(
		testUser.Username, testUser.Password, testUser.Email, testUser.Image, testUser.Id).WillReturnResult(
		sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := ud.Update(testUser.Id, &testUser); err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserDatabase_GetById(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	ud := NewUserDatabase(db, zap.NewExample())

	mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
		testUser.Id).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "username", "password", "email", "image"}).AddRow(
		testUser.Id, testUser.Username, testUser.Password, testUser.Email, testUser.Image))

	if usr, err := ud.GetById(testUser.Id); err != nil {
		t.Fatalf("unexpected error %s", err)
	} else {
		assert.Equal(t, testUser, *usr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserDatabase_GetByName(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	ud := NewUserDatabase(db, zap.NewExample())

	mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE \(username = (.*)\) ORDER BY "kinopoisk"."users"."id" ASC LIMIT 1`).WithArgs(
		testUser.Username).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "username", "password", "email", "image"}).AddRow(
		testUser.Id, testUser.Username, testUser.Password, testUser.Email, testUser.Image))

	if usr, err := ud.GetByName(testUser.Username); err != nil {
		t.Fatalf("unexpected error %s", err)
	} else {
		assert.Equal(t, testUser, *usr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
