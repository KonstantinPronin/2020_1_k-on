package server

import (
	"github.com/go-park-mail-ru/2020_1_k-on/internal/middleware"
	repository2 "github.com/go-park-mail-ru/2020_1_k-on/internal/session/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/user/delivery/http"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/user/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/user/usecase"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func NewServer(e *echo.Echo, db *gorm.DB) {
	sessions := repository2.NewSessionStorage()
	users := repository.NewUserDatabase(db)
	//users := repository.NewPostgresForUser(cp)
	auth := middleware.NewAuth(sessions)
	user := usecase.NewUser(sessions, users)

	http.NewUserHandler(e, user, auth)
}
