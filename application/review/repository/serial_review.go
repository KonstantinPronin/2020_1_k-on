package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type SerialReviewDatabase struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewSerialReviewDatabase(conn *gorm.DB, logger *zap.Logger) review.Repository {
	return &SerialReviewDatabase{
		conn:   conn,
		logger: logger,
	}
}

func (r *SerialReviewDatabase) Add(review *models.Review) error {
	return r.conn.Table("kinopoisk.serial_reviews").Create(review).Error
}

func (r *SerialReviewDatabase) GetByProductId(id uint) ([]models.Review, error) {
	var reviews []models.Review

	err := r.conn.Table("kinopoisk.serial_reviews").Where("product_id = ?", id).Find(&reviews).Error

	if err != nil {
		return nil, err
	}
	return reviews, nil
}
