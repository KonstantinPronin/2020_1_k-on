package models

type Film struct {
	ID              uint    `json:"id" gorm:"primary_key"`
	MainGenre       string  `json:"mainGenre" gorm:"column:maingenre"`
	RussianName     string  `json:"russianName" gorm:"column:russianname"`
	EnglishName     string  `json:"englishName" gorm:"column:englishname"`
	TrailerLink     string  `json:"trailerLink" gorm:"column:trailerlink"`
	Rating          float64 `json:"rating"`
	ImdbRating      float64 `json:"imdbRating" gorm:"column:imdbrating"`
	TotalVotes      int     `json:"totalVotes" gorm:"column:totalvotes"` //всего голосов
	SumVotes        int     `json:"-" gorm:"column:sumvotes"`            //сумма голосов,нужна только бэку
	Description     string  `json:"description"`
	Image           string  `json:"image,omitempty"`
	BackgroundImage string  `json:"backgroundImage,omitempty" gorm:"column:backgroundimage"`
	Country         string  `json:"country"`
	Year            int     `json:"year"`
	AgeLimit        int     `json:"ageLimit,omitempty" gorm:"column:agelimit"`
}

//easyjson:json
type Films []Film

func (f *Film) TableName() string {
	return "films"
}

type ListFilm struct {
	ID          uint    `json:"id" gorm:"primary_key"`
	RussianName string  `json:"russianName" gorm:"column:russianname"`
	Image       string  `json:"image,omitempty"`
	Country     string  `json:"country"`
	Year        int     `json:"year"`
	AgeLimit    int     `json:"ageLimit,omitempty" gorm:"column:agelimit"`
	Rating      float64 `json:"rating"`
}

//easyjson:json
type ListsFilm []ListFilm

func FormatFilm(film Film) (list ListFilm) {
	list.ID = film.ID
	list.RussianName = film.RussianName
	list.Image = film.Image
	list.Country = film.Country
	list.Year = film.Year
	list.AgeLimit = film.AgeLimit
	list.Rating = film.Rating
	return list
}

func (lists *ListsFilm) Convert(films Films) ListsFilm {
	var lf ListsFilm
	for _, f := range films {
		lf = append(lf, FormatFilm(f))
	}
	return lf
}

//func (film *Film) SetRating(rating float64) {
//	film.Rating = rating
//}
