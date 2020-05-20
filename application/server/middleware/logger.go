package middleware

import (
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) Logger {
	return Logger{logger: logger}
}

func (l *Logger) Log(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		l.logger.Info(ctx.Request().URL.String(),
			zap.String("method", ctx.Request().Method),
			zap.String("host", ctx.Request().Host),
		)

		err := next(ctx)

		if err != nil {
			var status int
			switch err.(type) {
			case *echo.HTTPError:
				httpErr := err
				status = httpErr.(*echo.HTTPError).Code
			default:
				status = ctx.Response().Status
			}

			l.logger.Debug(ctx.Request().URL.String(),
				zap.Int("Status", status),
				zap.Int64("size", ctx.Response().Size),
			)
			return err
		}
		return nil
	}
}
