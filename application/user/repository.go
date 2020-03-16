package user

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type Repository interface {
	Add(user *models.User) (err error)
	Update(id int64, upUser *models.User) error
	GetById(id int64) (user *models.User, err error)
	GetByName(login string) (user *models.User, err error)
	Contains(login string) (bool, error)
}
