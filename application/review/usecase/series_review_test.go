package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review/mocks"
	mock_series "github.com/go-park-mail-ru/2020_1_k-on/application/series/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testSeries = models.Series{ID: 1}

func beforeTest(t *testing.T) (*mock_series.MockRepository, *mocks.MockRepository, review.UseCase) {
	ctrl := gomock.NewController(t)
	s := mock_series.NewMockRepository(ctrl)
	r := mocks.NewMockRepository(ctrl)

	uc := NewSeriesReview(r, s)
	return s, r, uc
}

func TestSeriesReview_Add_WrongId(t *testing.T) {
	_, _, uc := beforeTest(t)

	testReview.UserId = 0
	err := uc.Add(&testReview)
	testReview.UserId = 1

	assert.NotNil(t, err)
}

func TestSeriesReview_Add(t *testing.T) {
	s, r, uc := beforeTest(t)

	s.EXPECT().GetSeriesByID(testReview.ProductId).Return(testSeries, true)
	r.EXPECT().Add(&testReview).Return(nil)

	err := uc.Add(&testReview)

	assert.Nil(t, err)
}

func TestSeriesReview_GetByProductId(t *testing.T) {
	s, r, uc := beforeTest(t)

	s.EXPECT().GetSeriesByID(testReview.ProductId).Return(testSeries, true)
	r.EXPECT().GetByProductId(testReview.ProductId).Return([]models.Review{testReview}, nil)

	res, err := uc.GetByProductId(testReview.ProductId)

	assert.Nil(t, err)
	assert.Equal(t, res, []models.Review{testReview})
}

func TestSeriesReview_GetReview(t *testing.T) {
	_, r, uc := beforeTest(t)

	r.EXPECT().GetReview(testReview.ProductId, testReview.UserId).Return(&testReview, nil)

	res, err := uc.GetReview(testReview.ProductId, testReview.UserId)

	assert.Nil(t, err)
	assert.Equal(t, testReview, *res)
}

func TestSeriesReview_GetReview_WrongId(t *testing.T) {
	_, _, uc := beforeFilmTest(t)

	_, err := uc.GetReview(0, testReview.UserId)

	assert.NotNil(t, err)
}
