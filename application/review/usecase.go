package review

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type UseCase interface {
	Add(review *models.Review) error
	GetByFilmId(id uint) ([]models.Review, error)
}
