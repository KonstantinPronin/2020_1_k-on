package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type SeriesReviewDatabase struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewSeriesReviewDatabase(conn *gorm.DB, logger *zap.Logger) review.Repository {
	return &SeriesReviewDatabase{
		conn:   conn,
		logger: logger,
	}
}

func (r *SeriesReviewDatabase) Add(review *models.Review) error {
	return r.conn.Table("kinopoisk.series_reviews").Create(review).Error
}

func (r *SeriesReviewDatabase) GetByProductId(id uint) ([]models.Review, error) {
	var reviews []models.Review

	rows, err := r.conn.Table("kinopoisk.series_reviews rev").
		Select("rev.id, rev.rating, rev.body, rev.user_id, rev.product_id, usr.username, usr.image").
		Joins("inner join kinopoisk.users usr on usr.id = rev.user_id").
		Where("rev.product_id = ?", id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *SeriesReviewDatabase) GetReview(productId uint, userId uint) (*models.Review, error) {
	rev := new(models.Review)

	rows, err := r.conn.Table("kinopoisk.series_reviews rev").
		Select("rev.id, rev.rating, rev.body, rev.user_id, rev.product_id, usr.username, usr.image").
		Joins("inner join kinopoisk.users usr on usr.id = rev.user_id").
		Where("rev.product_id = ? and rev.user_id = ?", productId, userId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&rev.Id, &rev.Rating, &rev.Body, &rev.UserId, &rev.ProductId, &rev.Usr.Username, &rev.Usr.Image)
	} else {
		return nil, errors.NewNotFoundError("no user reviews")
	}
	if err != nil {
		return nil, err
	}

	return rev, nil
}
