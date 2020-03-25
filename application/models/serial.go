package models

type Serial struct {
	ID              uint    `json:"ID" gorm:"primary_key"`
	MainGenre       string  `json:"maingenre" gorm:"column:maingenre"` //русское
	RussianName     string  `json:"russianname" gorm:"column:russianname"`
	EnglishName     string  `json:"englishname" gorm:"column:englishname"`
	TrailerLink     string  `json:"trailerlink" gorm:"column:trailerlink"`
	Rating          float64 `json:"rating"`
	ImdbRating      float64 `json:"imdbrating" gorm:"column:imdbrating"`
	TotalVotes      int     `json:"totalvotes" gorm:"column:totalvotes"` //всего голосов
	SumVotes        int     `json:"sumvotes" gorm:"column:sumvotes"`     //сумма голосов,нужна только бэку
	Description     string  `json:"description"`
	Image           string  `json:"image,omitempty"`
	BackgroundImage string  `json:"backgroundimage,omitempty" gorm:"column:backgroundimage"`
	Country         string  `json:"country"`
	YearFirst       int     `json:"yearfirst" gorm:"column:yearfirst"` //начало
	YearLast        int     `json:"yearlast" gorm:"column:yearlast"`   //0 - "не закончился"
	AgeLimit        int     `json:"agelimit,omitempty" gorm:"column:agelimit"`
}

//easyjson:json
type Serials []Serial

func (s *Serial) TableName() string {
	return "serials"
}

type ListSerial struct {
	ID          uint    `json:"ID" gorm:"primary_key"`
	RussianName string  `json:"russianname" gorm:"column:russianname"`
	Image       string  `json:"image,omitempty"`
	Country     string  `json:"country"`
	YearFirst   int     `json:"yearfirst" gorm:"column:yearfirst"` //начало
	YearLast    int     `json:"yearlast" gorm:"column:yearlast"`   //0-  "не закончился"
	AgeLimit    int     `json:"agelimit,omitempty" gorm:"column:agelimit"`
	Rating      float64 `json:"rating"`
}

//easyjson:json
type ListsSerials []ListSerial

func FormatSerial(serial Serial) (list ListSerial) {
	list.ID = serial.ID
	list.RussianName = serial.RussianName
	list.Image = serial.Image
	list.Country = serial.Country
	list.YearFirst = serial.YearFirst
	list.YearLast = serial.YearLast
	list.AgeLimit = serial.AgeLimit
	list.Rating = serial.Rating
	return list
}

func (lists *ListsSerials) Convert(serials Serials) ListsSerials {
	var lf ListsSerials
	for _, s := range serials {
		lf = append(lf, FormatSerial(s))
	}
	return lf
}
