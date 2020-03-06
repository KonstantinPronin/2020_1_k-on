package models

type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}

func Generate(status int, body interface{}) Response {
	return Response{
		Status: status,
		Body:   body,
	}
}

//easyjson:json
type Responses []Response
