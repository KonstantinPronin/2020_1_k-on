package http

import (
	"bytes"
	"fmt"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	mocks2 "github.com/go-park-mail-ru/2020_1_k-on/application/server/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/application/session"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"testing"
	"time"
)

var (
	testUser = models.User{
		Username: "test",
		Password: "test",
		Email:    "test@example.com",
	}
	cookie = http.Cookie{
		Name:    session.CookieName,
		Value:   sessionId,
		Path:    "/",
		Expires: time.Now().Add(session.CookieDuration),
	}
)

const (
	errMsg    = "test error"
	sessionId = "test"
	ok        = 200
)

func beforeTest(t *testing.T) (*UserHandler, *mocks2.MockContext, *mocks.MockUseCase, *mocks2.MockResponseWriter) {
	ctrl := gomock.NewController(t)
	w := mocks2.NewMockResponseWriter(ctrl)
	ctx := mocks2.NewMockContext(ctrl)
	uc := mocks.NewMockUseCase(ctrl)
	response := echo.NewResponse(w, echo.New())

	ctx.EXPECT().Response().Return(response).AnyTimes()
	return &UserHandler{useCase: uc, logger: zap.NewExample(), sanitizer: bluemonday.UGCPolicy()}, ctx, uc, w
}

func TestUserHandler_Login(t *testing.T) {
	uh, ctx, uc, w := beforeTest(t)

	body, err := easyjson.Marshal(testUser)
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}
	request, err := http.NewRequest("", "", bytes.NewReader(body))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request)
	uc.EXPECT().Login(testUser.Username, testUser.Password).Return(sessionId, "", nil)
	ctx.EXPECT().SetCookie(gomock.Any()).Do(func(arg *http.Cookie) {
		assert.Equal(t, sessionId, arg.Value)
	})
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Header().Return(http.Header{})
	w.EXPECT().Write(gomock.Any())

	err = uh.Login(ctx)
	assert.Nil(t, err)
}

func TestUserHandler_Logout_WithoutCookie(t *testing.T) {
	uh, ctx, _, w := beforeTest(t)

	ctx.EXPECT().Cookie(gomock.Any()).Return(nil, fmt.Errorf(errMsg))
	w.EXPECT().WriteHeader(http.StatusUnauthorized)
	w.EXPECT().Write(gomock.Any())

	expErr := echo.NewHTTPError(http.StatusUnauthorized, "no cookie")
	err := uh.Logout(ctx)
	assert.Equal(t, expErr, err)
}

func TestUserHandler_Logout(t *testing.T) {
	uh, ctx, uc, w := beforeTest(t)

	ctx.EXPECT().Cookie(gomock.Any()).Return(&cookie, nil)
	uc.EXPECT().Logout(gomock.Any()).Return(nil)
	ctx.EXPECT().SetCookie(gomock.Any()).Do(func(arg *http.Cookie) {
		assert.Equal(t, sessionId, arg.Value)
		assert.True(t, arg.Expires.Before(time.Now()))
	})
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Write(gomock.Any())

	err := uh.Logout(ctx)
	assert.Nil(t, err)
}

func TestUserHandler_SignUp(t *testing.T) {
	uh, ctx, uc, w := beforeTest(t)

	body, err := easyjson.Marshal(testUser)
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}
	request, err := http.NewRequest("", "", bytes.NewReader(body))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	uc.EXPECT().Add(&testUser).Return(&testUser, nil)
	uc.EXPECT().Login(testUser.Username, testUser.Password).Return(sessionId, "", nil)
	ctx.EXPECT().SetCookie(gomock.Any()).Do(func(arg *http.Cookie) {
		assert.Equal(t, sessionId, arg.Value)
	})
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Header().Return(http.Header{})
	w.EXPECT().Write(gomock.Any())

	err = uh.SignUp(ctx)
	assert.Nil(t, err)
}

func TestUserHandler_Profile(t *testing.T) {
	uh, ctx, uc, w := beforeTest(t)
	id := uint(0)

	ctx.EXPECT().Get(gomock.Any()).Return(id)
	uc.EXPECT().Get(id).Return(&testUser, nil)
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Write(gomock.Any())

	err := uh.Profile(ctx)
	assert.Nil(t, err)
}

func TestUserHandler_Update(t *testing.T) {
	uh, ctx, uc, w := beforeTest(t)

	id := uint(0)
	body, err := easyjson.Marshal(testUser)
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}
	request, err := http.NewRequest("", "", bytes.NewReader(body))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request)
	ctx.EXPECT().Get(gomock.Any()).Return(id)
	uc.EXPECT().Update(&testUser).Return(nil)
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Write(gomock.Any())

	err = uh.Update(ctx)
	assert.Nil(t, err)
}
