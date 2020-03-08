package http

import (
	"github.com/go-park-mail-ru/2020_1_k-on/internal/middleware"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/models"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/session"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/user"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

type userHandler struct {
	useCase user.UseCase
}

func NewUserHandler(e *echo.Echo, us user.UseCase, auth middleware.Auth) {
	handler := userHandler{useCase: us}

	e.Use(middleware.ParseErrors)
	e.POST("/login", handler.Login, auth.AlreadyLoginErr)
	e.POST("/logout", handler.Logout, auth.GetSession)
	e.PUT("/signup", handler.Signup, auth.AlreadyLoginErr)
	e.GET("/user", handler.Profile, auth.GetSession)
	e.POST("/user", handler.Update, auth.GetSession)
}

func (uh *userHandler) Login(ctx echo.Context) error {
	usr := new(models.User)
	if err := ctx.Bind(usr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "request parser error")
	}

	sessionId, err := uh.useCase.Login(usr.Username, usr.Password)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:    session.CookieName,
		Value:   sessionId,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Hour),
	}

	ctx.SetCookie(cookie)
	return nil
}

func (uh *userHandler) Logout(ctx echo.Context) error {
	sessionId := ctx.Get(session.CookieName)
	return uh.useCase.Logout(sessionId.(string))
}

func (uh *userHandler) Signup(ctx echo.Context) error {
	usr := new(models.User)
	if err := ctx.Bind(usr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "request parser error")
	}

	usr, err := uh.useCase.Add(usr)
	if err != nil {
		return err
	}

	usr.Password = ""
	return ctx.JSON(http.StatusOK, usr)
}

func (uh *userHandler) Profile(ctx echo.Context) error {
	uid := ctx.Get(session.UserIdKey)

	usr, err := uh.useCase.Get(uid.(int64))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	usr.Password = ""
	return ctx.JSON(http.StatusOK, usr)
}

func (uh *userHandler) Update(ctx echo.Context) error {
	uid := ctx.Get(session.UserIdKey)
	usr := new(models.User)
	if err := ctx.Bind(usr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "request parser error")
	}

	usr.Id = uid.(int64)
	err := uh.useCase.Update(usr)
	if err != nil {
		return err
	}

	usr.Password = ""
	return ctx.JSON(http.StatusOK, usr)
}
