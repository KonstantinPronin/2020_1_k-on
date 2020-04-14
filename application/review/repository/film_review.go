package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type FilmReviewDatabase struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewFilmReviewDatabase(conn *gorm.DB, logger *zap.Logger) review.Repository {
	return &FilmReviewDatabase{
		conn:   conn,
		logger: logger,
	}
}

func (r *FilmReviewDatabase) Add(review *models.Review) error {
	return r.conn.Table("kinopoisk.film_reviews").Create(review).Error
}

func (r *FilmReviewDatabase) GetByProductId(id uint) ([]models.Review, error) {
	var reviews []models.Review

	rows, err := r.conn.Table("kinopoisk.film_reviews rev").
		Select("rev.id, rev.rating, rev.body, rev.user_id, rev.product_id, usr.username, usr.image").
		Joins("inner join kinopoisk.users usr on usr.id = rev.user_id").
		Where("rev.product_id = ?", id).Rows()
	if err != nil {
		return nil, err
	}

	rev := new(models.Review)
	usr := new(models.ListUser)
	for rows.Next() {
		err = rows.Scan(&rev.Id, &rev.Rating, &rev.Body, &rev.UserId, &rev.ProductId, &usr.Username, &usr.Image)
		if err != nil {
			return nil, err
		}

		rev.Usr = *usr
		reviews = append(reviews, *rev)
	}

	return reviews, nil
}

func (r *FilmReviewDatabase) GetReview(productId uint, userId uint) (*models.Review, error) {
	rev := new(models.Review)

	err := r.conn.Table("kinopoisk.film_reviews rev").
		Select("rev.id, rev.rating, rev.body, rev.user_id, rev.product_id, usr.username, usr.image").
		Joins("inner join kinopoisk.users usr on usr.id = rev.user_id").
		Where("rev.product_id = ? and rev.user_id = ?", productId, userId).Row().
		Scan(&rev.Id, &rev.Rating, &rev.Body, &rev.UserId, &rev.ProductId, &rev.Usr.Username, &rev.Usr.Image)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewNotFoundError(err.Error())
		}
		return nil, err
	}

	return rev, nil
}
