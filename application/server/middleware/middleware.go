package middleware

import (
	"encoding/json"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"log"
)

var logConf = []byte(`{
    "level": "debug",
    "encoding": "console",
    "outputPaths": ["stdout", "/tmp/logs"],
    "errorOutputPaths": ["stderr"],
    "encoderConfig": {
      "messageKey": "message",
          "callerKey": "caller",
      "callerEncoder": "short",
      "levelKey": "level",
      "levelEncoder": "capital",
      "timeKey": "time",
       "timeEncoder": "ISO8601"
    }
    }`)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		//log request example
		var cfg zap.Config
		if err := json.Unmarshal(logConf, &cfg); err != nil {
			log.Print("config problem in middleware")
		}
		zapLogger, err := cfg.Build()
		if err != nil {
			log.Print("logger build problem in middleware")
		}
		defer zapLogger.Sync()

		zapLogger.Info(ctx.Request().URL.String(),
			zap.String("method", ctx.Request().Method),
			zap.String("host", ctx.Request().Host),
		)
		err = next(ctx)
		if err != nil {
			//log response
			zapLogger.Debug(ctx.Request().URL.String(),
				zap.Int("Status", ctx.Response().Status),
				zap.Int64("size", ctx.Response().Size),
			)
			return err
		}
		return nil
	}
}

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
