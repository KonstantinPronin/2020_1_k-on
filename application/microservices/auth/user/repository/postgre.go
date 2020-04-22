package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/user"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type UserDatabase struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewUserDatabase(db *gorm.DB, logger *zap.Logger) user.Repository {
	return &UserDatabase{conn: db, logger: logger}
}

func (udb *UserDatabase) GetByName(login string) (usr *models.User, err error) {
	usr = new(models.User)
	err = udb.conn.Table("kinopoisk.users").Where("username = ?", login).First(usr).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.NewNotFoundError(err.Error())
	}
	return usr, err
}
