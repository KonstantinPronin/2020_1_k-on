package user

import (
	"github.com/go-park-mail-ru/2020_1_k-on/internal/models"
)

type UseCase interface {
	Login(login string, password string) (sessionId string, err error)
	Logout(sessionId string) error
	Add(usr *models.User) (*models.User, error)
	Get(id int64) (*models.User, error)
	Update(user *models.User) error
}
