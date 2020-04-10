package user

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type Repository interface {
	Add(user *models.User) (err error)
	Update(id uint, upUser *models.User) error
	GetById(id uint) (user *models.User, err error)
	GetByName(login string) (user *models.User, err error)
	Contains(login string) (bool, error)
	SetImage(id uint, image string) error
}
