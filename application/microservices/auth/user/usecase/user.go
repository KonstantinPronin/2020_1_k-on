package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/session"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/user"
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

func (uc *User) Login(login string, password string) (sessionId string, csrfToken string, err error) {
	if login == "" || password == "" {
		return "", "", errors.NewInvalidArgument("Empty login or password")
	}

	usr, err := uc.users.GetByName(login)
	if err != nil {
		return "", "", err
	}

	ok, err := crypto.CheckPassword(password, usr.Password)
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", errors.NewInvalidArgument("Wrong password")
	}

	sessionId = uuid.New().String()
	csrfToken = crypto.CreateToken(sessionId)
	err = uc.sessions.Add(sessionId, usr.Id)

	return sessionId, csrfToken, err
}

func (uc *User) Check(sessionId string) (uint, error) {
	return uc.sessions.GetUserId(sessionId)
}

func (uc *User) Logout(sessionId string) error {
	if sessionId == "" {
		return errors.NewInvalidArgument("Empty session id")
	}
	return uc.sessions.Delete(sessionId)
}
