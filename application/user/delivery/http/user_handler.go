package http

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	"github.com/go-park-mail-ru/2020_1_k-on/application/session"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type UserHandler struct {
	useCase user.UseCase
	logger  *zap.Logger
}

func NewUserHandler(e *echo.Echo, us user.UseCase, auth middleware.Auth, logger *zap.Logger) {
	handler := UserHandler{useCase: us, logger: logger}

	//e.Use(middleware.ParseErrors)
	e.POST("/login", handler.Login, auth.AlreadyLoginErr, middleware.ParseErrors)
	e.POST("/logout", handler.Logout, auth.GetSession, middleware.ParseErrors)
	e.PUT("/signup", handler.SignUp, auth.AlreadyLoginErr, middleware.ParseErrors)
	e.GET("/user", handler.Profile, auth.GetSession, middleware.ParseErrors)
	e.POST("/user", handler.Update, auth.GetSession, middleware.ParseErrors)
}

func (uh *UserHandler) Login(ctx echo.Context) error {
	usr := new(models.User)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, usr); err != nil {
		uh.logger.Error("request parser error")
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "request parser error")
	}

	sessionId, err := uh.useCase.Login(usr.Username, usr.Password)
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
	}

	uh.setCookie(ctx, sessionId)

	usr.Password = ""
	return middleware.WriteOkResponse(ctx, usr)
}

func (uh *UserHandler) Logout(ctx echo.Context) error {
	cookie, err := ctx.Cookie(session.CookieName)
	if err != nil {
		uh.logger.Warn("request without cookie")
		return middleware.WriteErrResponse(ctx, http.StatusUnauthorized, "no cookie")
	}

	if err := uh.useCase.Logout(cookie.Value); err != nil {
		return err
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	ctx.SetCookie(cookie)
	return middleware.WriteOkResponse(ctx, nil)
}

func (uh *UserHandler) SignUp(ctx echo.Context) error {
	usr := new(models.User)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, usr); err != nil {
		uh.logger.Error("request parser error")
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, "request parser error")
	}

	password := usr.Password
	usr, err := uh.useCase.Add(usr)
	if err != nil {
		return err
	}

	sessionId, err := uh.useCase.Login(usr.Username, password)
	if err != nil {
		return err
	}

	uh.setCookie(ctx, sessionId)

	usr.Password = ""
	return middleware.WriteOkResponse(ctx, usr)
}

func (uh *UserHandler) Profile(ctx echo.Context) error {
	uid := ctx.Get(session.UserIdKey)
	usr, err := uh.useCase.Get(uid.(int64))
	if err != nil {
		return err
	}

	usr.Password = ""
	return middleware.WriteOkResponse(ctx, usr)
}

func (uh *UserHandler) Update(ctx echo.Context) error {
	uid := ctx.Get(session.UserIdKey)
	usr := new(models.User)
	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, usr); err != nil {
		uh.logger.Error("request parser error")
		return echo.NewHTTPError(http.StatusBadRequest, "request parser error")
	}

	usr.Id = uid.(int64)
	err := uh.useCase.Update(usr)
	if err != nil {
		return err
	}

	usr.Password = ""
	return middleware.WriteOkResponse(ctx, usr)
}

func (uh *UserHandler) setCookie(ctx echo.Context, sessionId string) {
	cookie := &http.Cookie{
		Name:     session.CookieName,
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(session.CookieDuration),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}
	ctx.SetCookie(cookie)
}
