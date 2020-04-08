package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestSeriesReviewDatabase_Add(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	reviews := NewSeriesReviewDatabase(db, zap.NewExample())

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*)"series_reviews"`).WithArgs(
		testReview.Id, testReview.Rating, testReview.Body, testReview.UserId, testReview.ProductId).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	if err := reviews.Add(&testReview); err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSeriesReviewDatabase_GetByProductId(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	reviews := NewSeriesReviewDatabase(db, zap.NewExample())

	mock.ExpectQuery(`SELECT (.*)id, (.*)rating, (.*)body, (.*)user_id, (.*)product_id, (.*)username, (.*)image ` +
		`FROM (.*)series_reviews (.*) inner join (.*)users (.*) on (.*)id = (.*)user_id WHERE (.*)product_id (.*)`).
		WithArgs(testReview.ProductId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "rating", "body", "user_id", "product_id", "username", "image"}).
			AddRow(testReview.Id, testReview.Rating, testReview.Body,
				testReview.UserId, testReview.ProductId, testReview.Usr.Username, testReview.Usr.Image))

	rev, err := reviews.GetByProductId(testReview.ProductId)

	assert.Nil(t, err)
	assert.Equal(t, rev[0], testReview)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSeriesReviewDatabase_GetReview(t *testing.T) {
	db, mock := initMockDb(t)
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("error '%s' while closing resource", err)
		}
	}()
	reviews := NewSeriesReviewDatabase(db, zap.NewExample())

	mock.ExpectQuery(`SELECT (.*) FROM (.*)series_reviews(.*) WHERE (.*)product_id = (.*) and user_id = (.*)`).
		WithArgs(testReview.ProductId, testReview.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "rating", "body", "user_id", "product_id", "username", "image"}).
			AddRow(testReview.Id, testReview.Rating, testReview.Body,
				testReview.UserId, testReview.ProductId, testReview.Usr.Username, testReview.Usr.Image))

	rev, err := reviews.GetReview(testReview.ProductId, testReview.UserId)

	testReview.Usr.Username = ""
	assert.Nil(t, err)
	assert.Equal(t, testReview, *rev)
	testReview.Usr.Username = "test"

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
