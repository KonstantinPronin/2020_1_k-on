package middleware

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/client"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/constants"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

type Auth struct {
	rpcAuth *client.AuthClient
}

func NewAuth(rpcAuth *client.AuthClient) Auth {
	return Auth{rpcAuth: rpcAuth}
}

func (a *Auth) GetSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie(constants.CookieName)
		if err != nil || cookie == nil {
			return WriteErrResponse(ctx, http.StatusUnauthorized, "no cookie")
		}

		uid, err := a.rpcAuth.Check(cookie.Value)
		if err != nil {
			cookie.Expires = time.Now().AddDate(0, 0, -1)
			ctx.SetCookie(cookie)
			return WriteErrResponse(ctx, http.StatusUnauthorized, "session does not exist")
		}

		ctx.Set(constants.CookieName, cookie.Value)
		ctx.Set(constants.UserIdKey, uid)
		return next(ctx)
	}
}

func (a *Auth) AlreadyLoginErr(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, err := ctx.Cookie(constants.CookieName)
		if err != http.ErrNoCookie {
			return WriteErrResponse(ctx, http.StatusForbidden, "already login")
		}

		return next(ctx)
	}
}
