package review

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type Repository interface {
	Add(review *models.Review) error
	GetByProductId(id uint) ([]models.Review, error)
	GetReview(productId uint, userId uint) (*models.Review, error)
}
