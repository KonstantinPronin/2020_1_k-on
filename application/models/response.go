package models

import "github.com/labstack/echo"

type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}

func Generate(status int, body interface{}, ctx *echo.Context) Response {
	(*ctx).Response().Status = status
	return Response{
		Status: status,
		Body:   body,
	}
}

//easyjson:json
type Responses []Response
