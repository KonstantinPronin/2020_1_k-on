package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
)

type SerialReview struct {
	reviews review.Repository
	series  series.Repository
}

func NewSerialReview(r review.Repository, s series.Repository) review.UseCase {
	return &SerialReview{
		reviews: r,
		series:  s,
	}
}

func (r *SerialReview) Add(review *models.Review) error {
	if review.ProductId == 0 || review.UserId == 0 {
		return errors.NewInvalidArgument("empty film id or user id")
	}

	s, _ := r.series.GetSeriesByID(review.ProductId)
	if s.ID != review.ProductId {
		return errors.NewInvalidArgument("wrong film id")
	}

	return r.reviews.Add(review)
}

func (r *SerialReview) GetByProductId(id uint) ([]models.Review, error) {
	s, _ := r.series.GetSeriesByID(id)
	if s.ID != id {
		return nil, errors.NewInvalidArgument("wrong film id")
	}

	return r.reviews.GetByProductId(id)
}
