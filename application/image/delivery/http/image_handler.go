package http

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/image"
	"github.com/go-park-mail-ru/2020_1_k-on/application/server/middleware"
	"github.com/go-park-mail-ru/2020_1_k-on/application/session"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
)

type ImageHandler struct {
	user   user.UseCase
	image  image.UseCase
	logger *zap.Logger
}

func NewUserHandler(e *echo.Echo, image image.UseCase, user user.UseCase, auth middleware.Auth, logger *zap.Logger) {
	handler := ImageHandler{
		user:   user,
		image:  image,
		logger: logger,
	}

	e.GET("/:path", handler.Get, middleware.ParseErrors)
	e.POST("/user/image", handler.AddUserImage, auth.GetSession, middleware.ParseErrors)
}

func (handler *ImageHandler) AddUserImage(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return middleware.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	path, err := handler.image.Save(src)
	if err != nil {
		return err
	}

	userId := ctx.Get(session.UserIdKey).(uint)
	err = handler.user.SetImage(userId, path)
	if err != nil {
		return err
	}

	return middleware.WriteOkResponse(ctx, path)
}

func (handler *ImageHandler) Get(ctx echo.Context) error {
	path := ctx.Param("path")
	path = filepath.Join(image.Dir, path)
	return ctx.File(path)
}
