package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
)

type SeriesReview struct {
	reviews review.Repository
	series  series.Repository
}

func NewSeriesReview(r review.Repository, s series.Repository) review.UseCase {
	return &SeriesReview{
		reviews: r,
		series:  s,
	}
}

func (r *SeriesReview) Add(review *models.Review) error {
	if review.ProductId == 0 || review.UserId == 0 {
		return errors.NewInvalidArgument("empty film id or user id")
	}

	s, _ := r.series.GetSeriesByID(review.ProductId)
	if s.ID != review.ProductId {
		return errors.NewInvalidArgument("wrong film id")
	}

	return r.reviews.Add(review)
}

func (r *SeriesReview) GetByProductId(id uint) ([]models.Review, error) {
	s, _ := r.series.GetSeriesByID(id)
	if s.ID != id {
		return nil, errors.NewInvalidArgument("wrong film id")
	}

	return r.reviews.GetByProductId(id)
}

func (r *SeriesReview) GetReview(productId uint, userId uint) (*models.Review, error) {
	if productId == 0 || userId == 0 {
		return nil, errors.NewInvalidArgument("wrong id")
	}

	return r.reviews.GetReview(productId, userId)
}
