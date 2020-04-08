package models

type Collection struct {
	Name string      `json:"name"`
	List interface{} `json:"list"`
}
