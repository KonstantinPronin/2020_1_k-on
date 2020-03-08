package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/internal/models"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/session"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/user"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/google/uuid"
)

type User struct {
	sessions session.Repository
	users    user.Repository
}

func NewUser(s session.Repository, u user.Repository) User {
	return User{sessions: s, users: u}
}

func (uc *User) Login(login string, password string) (sessionId string, err error) {
	if login == "" || password == "" {
		return "", errors.NewInvalidArgument("Empty login or password")
	}

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

	id, err := uc.users.Add(usr)
	if err != nil {
		return nil, err
	}

	usr.Id = id
	return usr, nil
}

func (uc *User) Get(id int64) (*models.User, error) {
	return uc.users.GetById(id)
}

func (uc *User) Update(usr *models.User) error {
	if usr.Id == -1 {
		return errors.NewInvalidArgument("User does not exist")
	}

	return uc.users.Update(usr.Id, usr)
}
