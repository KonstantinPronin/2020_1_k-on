package http

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/client"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/constants"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	vkc "golang.org/x/oauth2/vk"
	"net/http"
	"time"
)

type UserHandler struct {
	rpcAuth   client.IAuthClient
	useCase   user.UseCase
	logger    *zap.Logger
	sanitizer *bluemonday.Policy
	oauthConf oauth2.Config
}

func NewUserHandler(e *echo.Echo,
	rpcAuth client.IAuthClient,
	uc user.UseCase,
	auth middleware.Auth,
	logger *zap.Logger,
	sanitizer *bluemonday.Policy) error {

	conf, err := uc.GetOauthConfig()
	if err != nil {
		return err
	}
	oauth := oauth2.Config{
		ClientID:     conf.ClientId,
		ClientSecret: conf.ClientSecret,
		Endpoint:     vkc.Endpoint,
		RedirectURL:  conf.RedirectUrl,
		Scopes:       []string{"email"},
	}
	handler := UserHandler{rpcAuth: rpcAuth, useCase: uc, logger: logger, sanitizer: sanitizer, oauthConf: oauth}

	e.POST("/login", handler.Login, auth.AlreadyLoginErr, middleware.ParseErrors)
	e.DELETE("/logout", handler.Logout, auth.GetSession, middleware.ParseErrors)
	e.POST("/signup", handler.SignUp, auth.AlreadyLoginErr, middleware.ParseErrors)
	e.GET("/user", handler.Profile, auth.GetSession, middleware.ParseErrors)
	e.PUT("/user", handler.Update, auth.GetSession, middleware.ParseErrors, middleware.CSRF)
	e.GET("/vk", handler.GetVkRedirect, middleware.ParseErrors)
	e.GET("/oauth", handler.Oauth, middleware.ParseErrors)

	return nil
}

func (uh *UserHandler) Login(ctx echo.Context) error {
	usr := new(models.User)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, usr); err != nil {
		uh.logger.Error("request parser error")
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "request parser error")
	}
	uh.sanitize(usr)

	sessionId, token, err := uh.rpcAuth.Login(usr.Username, usr.Password)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
	}

	uh.setCookie(ctx, sessionId)
	uh.setCsrfToken(ctx, token)

	usr.Password = ""
	return middleware.WriteOkResponse(ctx, usr)
}

func (uh *UserHandler) Logout(ctx echo.Context) error {
	cookie, err := ctx.Cookie(constants.CookieName)
	if err != nil {
		uh.logger.Warn("request without cookie")
		return middleware.WriteErrResponse(ctx, http.StatusUnauthorized, "no cookie")
	}

	if err := uh.rpcAuth.Logout(cookie.Value); err != nil {
		return err
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	cookie.Path = "/"
	cookie.HttpOnly = true
	ctx.SetCookie(cookie)
	return middleware.WriteOkResponse(ctx, nil)
}

func (uh *UserHandler) SignUp(ctx echo.Context) error {
	usr := new(models.User)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, usr); err != nil {
		uh.logger.Error("request parser error")
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "request parser error")
	}
	uh.sanitize(usr)

	password := usr.Password
	usr, err := uh.useCase.Add(usr)
	if err != nil {
		return err
	}

	sessionId, token, err := uh.rpcAuth.Login(usr.Username, password)
	if err != nil {
		return err
	}

	uh.setCookie(ctx, sessionId)
	uh.setCsrfToken(ctx, token)

	usr.Password = ""
	return middleware.WriteOkResponse(ctx, usr)
}

func (uh *UserHandler) Profile(ctx echo.Context) error {
	uid := ctx.Get(constants.UserIdKey)
	usr, err := uh.useCase.Get(uid.(uint))
	if err != nil {
		return err
	}

	usr.Password = ""
	return middleware.WriteOkResponse(ctx, usr)
}

func (uh *UserHandler) Update(ctx echo.Context) error {
	uid := ctx.Get(constants.UserIdKey)
	usr := new(models.User)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, usr); err != nil {
		uh.logger.Error("request parser error")
		return echo.NewHTTPError(http.StatusBadRequest, "request parser error")
	}
	uh.sanitize(usr)

	usr.Id = uid.(uint)
	err := uh.useCase.Update(usr)
	if err != nil {
		return err
	}

	usr.Password = ""
	return middleware.WriteOkResponse(ctx, usr)
}

func (uh *UserHandler) GetVkRedirect(ctx echo.Context) error {
	url := uh.oauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return middleware.WriteOkResponse(ctx, url)
}

func (uh *UserHandler) Oauth(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	authToken, err := uh.oauthConf.Exchange(context.Background(), code)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusInternalServerError, "can not get auth token from vk")
	}

	vkUser := new(models.VkUser)
	vkUser.Id = int64(authToken.Extra("user_id").(float64))
	vkUser.Email = authToken.Extra("email").(string)

	usr, err := uh.useCase.Oauth(vkUser)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusInternalServerError, "can not create vk user")
	}

	sessionId, csrfToken, err := uh.rpcAuth.Login(usr.Username, usr.Password)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
	}

	uh.setCookie(ctx, sessionId)
	uh.setCsrfToken(ctx, csrfToken)

	usr.Password = ""
	return middleware.WriteOkResponse(ctx, usr)
}

func (uh *UserHandler) setCookie(ctx echo.Context, sessionId string) {
	cookie := &http.Cookie{
		Name:    constants.CookieName,
		Value:   sessionId,
		Path:    "/",
		Expires: time.Now().Add(constants.CookieDuration),
		//SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}
	ctx.SetCookie(cookie)
}

func (uh *UserHandler) setCsrfToken(ctx echo.Context, token string) {
	cookie := &http.Cookie{
		Name:    constants.CSRFHeader,
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
		//SameSite: http.SameSiteStrictMode,
	}
	ctx.SetCookie(cookie)
}

func (uh *UserHandler) sanitize(usr *models.User) {
	usr.Username = uh.sanitizer.Sanitize(usr.Username)
	usr.Password = uh.sanitizer.Sanitize(usr.Password)
	usr.Email = uh.sanitizer.Sanitize(usr.Email)
	usr.Image = uh.sanitizer.Sanitize(usr.Image)
}
