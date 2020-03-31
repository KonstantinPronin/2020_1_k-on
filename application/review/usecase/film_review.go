package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/film"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
)

type FilmReview struct {
	reviews review.Repository
	films   film.Repository
}

func NewFilmReview(r review.Repository, f film.Repository) review.UseCase {
	return &FilmReview{
		reviews: r,
		films:   f,
	}
}

func (r *FilmReview) Add(review *models.Review) error {
	if review.ProductId == 0 || review.UserId == 0 {
		return errors.NewInvalidArgument("empty film id or user id")
	}

	f, _ := r.films.GetById(review.ProductId)
	if f.ID != review.ProductId {
		return errors.NewInvalidArgument("wrong film id")
	}

	return r.reviews.Add(review)
}

func (r *FilmReview) GetByProductId(id uint) ([]models.Review, error) {
	f, _ := r.films.GetById(id)
	if f.ID != id {
		return nil, errors.NewInvalidArgument("wrong film id")
	}

	return r.reviews.GetByProductId(id)
}
