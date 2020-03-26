package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type ReviewDatabase struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewReviewDatabase(conn *gorm.DB, logger *zap.Logger) review.Repository {
	return &ReviewDatabase{
		conn:   conn,
		logger: logger,
	}
}

func (r *ReviewDatabase) Add(review *models.Review) error {
	return r.conn.Table("kinopoisk.reviews").Create(review).Error
}

func (r *ReviewDatabase) GetByFilmId(id uint) ([]models.Review, error) {
	var reviews []models.Review

	err := r.conn.Table("kinopoisk.reviews").Where("filmId = ?", id).Find(&reviews).Error

	if err != nil {
		return nil, err
	}
	return reviews, nil
}
