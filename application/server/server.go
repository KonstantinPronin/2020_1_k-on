package server

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/middleware"
	repository2 "github.com/go-park-mail-ru/2020_1_k-on/application/session/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user/delivery/http"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/user/usecase"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func NewServer(e *echo.Echo, db *gorm.DB, rd *redis.Client, logger *zap.Logger) {
	sessions := repository2.NewSessionDatabase(rd, logger)
	users := repository.NewUserDatabase(db, logger)
	auth := middleware.NewAuth(sessions)
	user := usecase.NewUser(sessions, users, logger)

	http.NewUserHandler(e, user, auth, logger)
}
