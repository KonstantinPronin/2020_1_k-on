package usecase

import (
	sessionmock "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/session/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user"
	usermock "github.com/go-park-mail-ru/2020_1_k-on/application/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

var testUser = models.User{
	Id:       1,
	Username: "test",
	Password: "test",
}

func beforeTest(t *testing.T) (*sessionmock.MockRepository, *usermock.MockRepository, user.UseCase) {
	ctrl := gomock.NewController(t)
	sessions := sessionmock.NewMockRepository(ctrl)
	users := usermock.NewMockRepository(ctrl)
	usecase := NewUser(users, zap.NewExample())

	return sessions, users, usecase
}

func TestUser_Add_EmptyInput(t *testing.T) {
	_, _, us := beforeTest(t)

	_, err := us.Add(&models.User{
		Username: "",
		Password: "",
	})

	assert.NotNil(t, err)
}

func TestUser_Add_UserAlreadyExist(t *testing.T) {
	_, users, us := beforeTest(t)

	users.EXPECT().Contains(testUser.Username).Return(true, nil)

	_, err := us.Add(&testUser)

	assert.NotNil(t, err)
}

func TestUser_Add_Success(t *testing.T) {
	_, users, us := beforeTest(t)

	users.EXPECT().Contains(testUser.Username).Return(false, nil)
	users.EXPECT().Add(&testUser).Return(nil)

	usr, err := us.Add(&testUser)

	assert.Nil(t, err)
	assert.Equal(t, testUser, *usr)
}

func TestUser_Get_Success(t *testing.T) {
	id := uint(0)
	_, users, us := beforeTest(t)

	users.EXPECT().GetById(id).Return(&testUser, nil)

	usr, err := us.Get(id)

	assert.Nil(t, err)
	assert.Equal(t, testUser, *usr)
}

func TestUser_Update_WrongUser(t *testing.T) {
	_, _, us := beforeTest(t)

	err := us.Update(&models.User{Id: 0})

	assert.NotNil(t, err)
}

func TestUser_Update_Success(t *testing.T) {
	_, users, us := beforeTest(t)

	users.EXPECT().Update(testUser.Id, &testUser).Return(nil)

	err := us.Update(&testUser)

	assert.Nil(t, err)
}
