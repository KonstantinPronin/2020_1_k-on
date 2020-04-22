package middleware

import (
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/constants"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/crypto"
	"github.com/labstack/echo"
	"net/http"
)

func CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("Content-Type", "application/json; charset=utf8")
		ctx.Response().Header().Set("Access-Control-Allow-Origin", ctx.Request().Header.Get("Origin"))
		//if ctx.Request().Method == "OPTIONS" {
		//	ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		//	ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//	ctx.Response().Header().Set("Access-Control-Allow-Origin", ctx.Request().Header.Get("Origin"))
		//	ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		//	return ctx.NoContent(200)
		//}
		return next(ctx)
	}
}

func CSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId := ctx.Get(constants.CookieName)
		token := ctx.Request().Header.Get(constants.CSRFHeader)
		if sessionId == "" || token == "" || !crypto.CheckToken(sessionId.(string), token) {
			return WriteErrResponse(ctx, http.StatusForbidden, "wrong csrf token")
		}

		return next(ctx)
	}
}
