package user

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type Repository interface {
	GetByName(login string) (user *models.User, err error)
}
