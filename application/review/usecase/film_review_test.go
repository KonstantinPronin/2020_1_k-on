package usecase

import (
	mock_film "github.com/go-park-mail-ru/2020_1_k-on/application/film/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testReview = models.Review{
	Rating:    10,
	Body:      "test",
	UserId:    1,
	ProductId: 1,
}

var testFilm = models.Film{ID: 1}

func beforeFilmTest(t *testing.T) (*mock_film.MockRepository, *mocks.MockRepository, review.UseCase) {
	ctrl := gomock.NewController(t)
	f := mock_film.NewMockRepository(ctrl)
	r := mocks.NewMockRepository(ctrl)

	uc := NewFilmReview(r, f)
	return f, r, uc
}

func TestFilmReview_Add_WrongId(t *testing.T) {
	_, _, uc := beforeFilmTest(t)

	testReview.UserId = 0
	err := uc.Add(&testReview)
	testReview.UserId = 1

	assert.NotNil(t, err)
}

func TestFilmReview_Add(t *testing.T) {
	f, r, uc := beforeFilmTest(t)

	f.EXPECT().GetById(testReview.ProductId).Return(&testFilm, true)
	r.EXPECT().Add(&testReview).Return(nil)

	err := uc.Add(&testReview)

	assert.Nil(t, err)
}

func TestFilmReview_GetByProductId(t *testing.T) {
	f, r, uc := beforeFilmTest(t)

	f.EXPECT().GetById(testReview.ProductId).Return(&testFilm, true)
	r.EXPECT().GetByProductId(testReview.ProductId).Return([]models.Review{testReview}, nil)

	res, err := uc.GetByProductId(testReview.ProductId)

	assert.Nil(t, err)
	assert.Equal(t, res, []models.Review{testReview})
}
