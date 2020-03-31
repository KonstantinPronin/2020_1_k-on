package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	sessionmock "github.com/go-park-mail-ru/2020_1_k-on/application/session/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user"
	usermock "github.com/go-park-mail-ru/2020_1_k-on/application/user/mocks"
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
	hash := sha256.Sum256([]byte(storedUser.Password))
	storedUser.Password = hex.EncodeToString(hash[:])

	users.EXPECT().GetByName(gomock.Eq(testUser.Username)).Return(&storedUser, nil)
	sessions.EXPECT().Add(gomock.Any(), testUser.Id).Return(nil)

	id, err := us.Login(testUser.Username, testUser.Password)

	assert.True(t, len(id) > 0)
	assert.Nil(t, err)
}

func TestUser_Login_EmptyLogin(t *testing.T) {
	_, _, us := beforeTest(t)

	id, err := us.Login("", testUser.Password)

	assert.True(t, len(id) == 0)
	assert.NotNil(t, err)
}

func TestUser_Login_UserNotFound(t *testing.T) {
	_, users, us := beforeTest(t)
	expectedErr := errors.NewNotFoundError("")

	users.EXPECT().GetByName(gomock.Eq(testUser.Username)).Return(nil, expectedErr)

	id, err := us.Login(testUser.Username, testUser.Password)

	assert.True(t, len(id) == 0)
	assert.Equal(t, expectedErr, err)
}

func TestUser_Login_WrongPassword(t *testing.T) {
	sessions, users, us := beforeTest(t)

	users.EXPECT().GetByName(gomock.Eq(testUser.Username)).Return(&testUser, nil)
	sessions.EXPECT().Add(gomock.Any(), testUser.Id).Return(nil)

	id, err := us.Login(testUser.Username, "wrong")

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
