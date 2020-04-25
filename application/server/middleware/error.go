package middleware

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"net/http"
)

func WriteOkResponse(ctx echo.Context, body interface{}) error {
	resp := models.Response{
		Status: http.StatusOK,
		Body:   body,
	}
	if _, err := easyjson.MarshalToWriter(resp, ctx.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func WriteErrResponse(ctx echo.Context, code int, message string) error {
	ctx.Response().Writer.WriteHeader(code)
	ctx.Response().Committed = true

	resp := models.Response{Status: code, Body: message}
	if _, err := easyjson.MarshalToWriter(resp, ctx.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echo.NewHTTPError(code, message)
}

func ParseErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)
		if err != nil {
			switch err.(type) {
			case *echo.HTTPError:
				return err
			case *errors.InvalidArgumentError:
				return WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
			case *errors.NotFoundError:
				return WriteErrResponse(ctx, http.StatusNotFound, err.Error())
			case *errors.ForbiddenError:
				return WriteErrResponse(ctx, http.StatusForbidden, err.Error())
			case *errors.DbInternalError:
				return WriteErrResponse(ctx, http.StatusInternalServerError, err.Error())
			default:
				return WriteErrResponse(ctx, http.StatusInternalServerError, err.Error())
			}
		}
		return nil
	}
}
