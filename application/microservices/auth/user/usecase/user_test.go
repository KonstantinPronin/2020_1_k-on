package usecase

import (
	sessionmock "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/session/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/user"
	usermock "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/user/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/crypto"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
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
	usecase := NewUser(sessions, users, zap.NewExample())

	return sessions, users, usecase
}

func TestUser_Login_Success(t *testing.T) {
	sessions, users, us := beforeTest(t)
	storedUser := testUser
	hash, err := crypto.HashPassword(testUser.Password)
	if err != nil {
		t.Error(err)
	}
	storedUser.Password = hash

	users.EXPECT().GetByName(gomock.Eq(testUser.Username)).Return(&storedUser, nil)
	sessions.EXPECT().Add(gomock.Any(), testUser.Id).Return(nil)

	id, _, err := us.Login(testUser.Username, testUser.Password)

	assert.True(t, len(id) > 0)
	assert.Nil(t, err)
}

func TestUser_Login_EmptyLogin(t *testing.T) {
	_, _, us := beforeTest(t)

	id, _, err := us.Login("", testUser.Password)

	assert.True(t, len(id) == 0)
	assert.NotNil(t, err)
}

func TestUser_Login_UserNotFound(t *testing.T) {
	_, users, us := beforeTest(t)
	expectedErr := errors.NewNotFoundError("")

	users.EXPECT().GetByName(gomock.Eq(testUser.Username)).Return(nil, expectedErr)

	id, _, err := us.Login(testUser.Username, testUser.Password)

	assert.True(t, len(id) == 0)
	assert.Equal(t, expectedErr, err)
}

func TestUser_Login_WrongPassword(t *testing.T) {
	sessions, users, us := beforeTest(t)

	users.EXPECT().GetByName(gomock.Eq(testUser.Username)).Return(&testUser, nil)
	sessions.EXPECT().Add(gomock.Any(), testUser.Id).Return(nil)

	id, _, err := us.Login(testUser.Username, "wrong")

	assert.True(t, len(id) == 0)
	assert.NotNil(t, err)
}

func TestUser_Logout_EmptySessionId(t *testing.T) {
	_, _, us := beforeTest(t)

	err := us.Logout("")

	assert.NotNil(t, err)
}

func TestUser_Logout_Success(t *testing.T) {
	sessionId := "test"
	sessions, _, us := beforeTest(t)

	sessions.EXPECT().Delete(sessionId).Return(nil)

	err := us.Logout(sessionId)

	assert.Nil(t, err)
}
