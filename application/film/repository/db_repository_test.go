package repository

import (
	"2020_1_k-on/application/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestPostgresForFilms_GetById(t *testing.T) {
	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_0")
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	db, err = gorm.Open("sqlmock", "sqlmock_db_0")
	if err != nil {
		t.Fatalf("Got an unexpected error: %s", err)

	}
	defer db.Close()

	var elemID uint = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name", "agelimit", "image"})
	//expect := *models.Film{elemID, "name", 10, "image"}
	expect := models.Film{elemID, "name", 10, "image"}
	//for _, item := range expect {
	rows = rows.AddRow(expect.ID, expect.Name, expect.AgeLimit, expect.Image)
	//}

	mock.ExpectQuery(`SELECT id,name,agelimit,image FROM "films" WHERE.*$`).
		//WithArgs(elemID).
		WillReturnRows(rows)

	repo := &PostgresForFilms{
		DB: db,
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
}
