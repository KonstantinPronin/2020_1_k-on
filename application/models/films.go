package models

type Film struct {
	ID          uint   `json:"ID"`
	Name        string `json:"name"`
	AgeLimit    int    `json:"agelimit,omitempty"`
	Image       string `json:"image,omitempty"`
	ImageBase64 string `json:"base,omitempty"`
}

//easyjson:json
type Films []Film
