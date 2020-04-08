package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/session"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user"
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

	password = uc.hashPassword(password)

	usr, err := uc.users.GetByName(login)
	if err != nil {
		return "", err
	}

	if usr.Password != password {
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

	usr.Password = uc.hashPassword(usr.Password)
	err := uc.users.Add(usr)
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

	usr.Password = uc.hashPassword(usr.Password)
	return uc.users.Update(usr.Id, usr)
}

func (uc *User) hashPassword(password string) string {
	if password == "" {
		return password
	}
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
