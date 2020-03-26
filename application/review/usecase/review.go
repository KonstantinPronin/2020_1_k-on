package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/film"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
)

type Review struct {
	reviews review.Repository
	films   film.Repository
}

func NewReview(r review.Repository, f film.Repository) review.UseCase {
	return &Review{
		reviews: r,
		films:   f,
	}
}

func (r *Review) Add(review *models.Review) error {
	if review.FilmId == 0 || review.UserId == 0 {
		return errors.NewInvalidArgument("empty film id or user id")
	}

	f, _ := r.films.GetById(review.FilmId)
	if f.ID != review.FilmId {
		return errors.NewInvalidArgument("wrong film id")
	}

	return r.reviews.Add(review)
}

func (r *Review) GetByFilmId(id uint) ([]models.Review, error) {
	f, _ := r.films.GetById(id)
	if f.ID != id {
		return nil, errors.NewInvalidArgument("wrong film id")
	}

	return r.reviews.GetByFilmId(id)
}
