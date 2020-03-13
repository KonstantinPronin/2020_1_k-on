package server

import (
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
		ctx.Response().Header().Set("Access-Control-Allow-Origin", ctx.Request().Header.Get("Origin"))
		ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Response().Header().Set("Content-Type", "application/json")
		//log request example
		zapLogger, _ := zap.NewProduction()
		defer zapLogger.Sync()
		zapLogger.Info(ctx.Request().URL.String(),
			zap.String("method", ctx.Request().Method),
			zap.String("host", ctx.Request().Host),
		)
		err := next(ctx)
		if err != nil {
			//log response
			zapLogger.Debug("",
				zap.Int("Status", ctx.Response().Status),
				zap.Int64("size", ctx.Response().Size),
			)
			return err
		}
		return nil
	}
}
