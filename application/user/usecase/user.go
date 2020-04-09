package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/session"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/crypto"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type User struct {
	sessions session.Repository
	users    user.Repository
	logger   *zap.Logger
}

func NewUser(s session.Repository, u user.Repository, logger *zap.Logger) user.UseCase {
	return &User{sessions: s, users: u, logger: logger}
}

func (uc *User) Login(login string, password string) (sessionId string, err error) {
	if login == "" || password == "" {
		return "", errors.NewInvalidArgument("Empty login or password")
	}

	usr, err := uc.users.GetByName(login)
	if err != nil {
		return "", err
	}

	ok, err := crypto.CheckPassword(password, usr.Password)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.NewInvalidArgument("Wrong password")
	}

	sessionId = uuid.New().String()
	err = uc.sessions.Add(sessionId, usr.Id)

	return sessionId, err
}

func (uc *User) Logout(sessionId string) error {
	if sessionId == "" {
		return errors.NewInvalidArgument("Empty session id")
	}
	return uc.sessions.Delete(sessionId)
}

func (uc *User) Add(usr *models.User) (*models.User, error) {
	if usr.Username == "" || usr.Password == "" {
		return nil, errors.NewInvalidArgument("Empty login or password")
	}

	if exist, err := uc.users.Contains(usr.Username); err != nil {
		return nil, err
	} else if exist {
		return nil, errors.NewInvalidArgument("User already exists")
	}

	hash, err := crypto.HashPassword(usr.Password)
	if err != nil {
		return nil, err
	}

	usr.Password = hash
	err = uc.users.Add(usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (uc *User) Get(id uint) (*models.User, error) {
	return uc.users.GetById(id)
}

func (uc *User) Update(usr *models.User) error {
	if usr.Id == 0 {
		uc.logger.Warn("User does not exist")
		return errors.NewInvalidArgument("User does not exist")
	}

	hash, err := crypto.HashPassword(usr.Password)
	if err != nil {
		return err
	}
	usr.Password = hash

	return uc.users.Update(usr.Id, usr)
}
