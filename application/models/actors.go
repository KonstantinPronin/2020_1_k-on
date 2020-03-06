package models

type Actor struct {
	ID         uint   `json:"ID"`
	Name       string `json:"name"`
	SecondName string `json:"secondname"`
	Age        int    `json:"age"`
	Image      string `json:"image,omitempty"`
}

//easyjson:json
type Actors []Actor
