package middleware

import (
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/labstack/echo"
	"net/http"
)

func ParseErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)
		if err != nil {
			switch err.(type) {
			case *echo.HTTPError:
				return err
			case *errors.InvalidArgumentError:
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			default:
				return echo.ErrInternalServerError
			}
		}
		return nil
	}
}
