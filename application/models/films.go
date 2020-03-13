package models

type Film struct {
	ID       uint   `json:"ID" gorm:"primary_key"`
	Name     string `json:"name"`
	AgeLimit int    `json:"agelimit,omitempty" gorm:"column:agelimit"`
	Image    string `json:"image,omitempty"`
}

//easyjson:json
type Films []Film

func (f *Film) TableName() string {
	return "films"
}
