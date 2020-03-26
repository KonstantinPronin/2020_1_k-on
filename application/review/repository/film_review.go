package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
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

	err := r.conn.Table("kinopoisk.film_reviews").Where("product_id = ?", id).Find(&reviews).Error

	if err != nil {
		return nil, err
	}
	return reviews, nil
}
