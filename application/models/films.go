package models

type Film struct {
	ID          uint    `json:"ID" gorm:"primary_key"`
	Type        string  `json:"type"`
	MainGenre   string  `json:"maingenre" gorm:"column:maingenre"`
	RussianName string  `json:"russianname" gorm:"column:russianname"`
	EnglishName string  `json:"englishname" gorm:"column:englishname"`
	Seasons     int     `json:"seasons"`
	TrailerLink string  `json:"trailerlink" gorm:"column:trailerlink"`
	Rating      float64 `json:"rating"`
	ImdbRating  float64 `json:"imdbrating" gorm:"column:imdbrating"`
	Description string  `json:"description"`
	Image       string  `json:"image,omitempty"`
	Country     string  `json:"country"`
	Year        int     `json:"year"`
	AgeLimit    int     `json:"agelimit,omitempty" gorm:"column:agelimit"`
}

//easyjson:json
type Films []Film

func (f *Film) TableName() string {
	return "films"
}

type ListFilm struct {
	ID          uint    `json:"ID" gorm:"primary_key"`
	Type        string  `json:"type"`
	EnglishName string  `json:"englishname" gorm:"column:englishname"`
	Image       string  `json:"image,omitempty"`
	MainGenre   string  `json:"maingenre" gorm:"column:maingenre"`
	AgeLimit    int     `json:"agelimit,omitempty" gorm:"column:agelimit"`
	Rating      float64 `json:"rating"`
}

//easyjson:json
type ListsFilm []ListFilm

func Format(film Film) (list ListFilm) {
	list.ID = film.ID
	list.Type = film.Type
	list.EnglishName = film.EnglishName
	list.Image = film.Image
	list.MainGenre = film.MainGenre
	list.AgeLimit = film.AgeLimit
	list.Rating = film.Rating
	return list
}

func (lists *ListsFilm) Convert(films Films) ListsFilm {
	var lf ListsFilm
	for _, f := range films {
		lf = append(lf, Format(f))
	}
	return lf
}
