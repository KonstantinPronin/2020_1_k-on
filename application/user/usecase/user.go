package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/crypto"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type User struct {
	users  user.Repository
	logger *zap.Logger
}

func NewUser(u user.Repository, logger *zap.Logger) user.UseCase {
	return &User{users: u, logger: logger}
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
	usr.Image = ""
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
	usr.Image = ""

	return uc.users.Update(usr.Id, usr)
}

func (uc *User) SetImage(id uint, image string) error {
	return uc.users.SetImage(id, image)
}

func (uc *User) Oauth(vkUser *models.VkUser) (*models.User, error) {
	usr, err := uc.users.GetUserFromVk(vkUser.Id)
	if err != nil {
		return nil, err
	}

	if usr != nil {
		password, hash, err := uc.generatePassword()
		if err != nil {
			return nil, err
		}

		usr.Password = hash
		err = uc.Update(usr)
		if err != nil {
			return nil, err
		}
		usr.Password = password

		return usr, nil
	}

	usr = new(models.User)
	usr.Username = vkUser.Email
	usr.Email = vkUser.Email
	usr.Image = ""

	password, hash, err := uc.generatePassword()
	if err != nil {
		return nil, err
	}

	usr.Password = hash
	err = uc.users.CreateUserFromVk(vkUser.Id, usr)
	if err != nil {
		return nil, err
	}
	usr.Password = password

	return usr, err
}

func (uc *User) GetOauthConfig() (*models.OauthConfig, error) {
	return uc.users.GetOauthConfig()
}

func (uc *User) generatePassword() (string, string, error) {
	password := uuid.New().String()
	hash, err := crypto.HashPassword(password)
	if err != nil {
		return "", "", err
	}

	return password, hash, err
}
