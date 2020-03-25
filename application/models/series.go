package models

//СЕРИАЛЫ!

type Series struct {
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
type SeriesArr []Series

func (s *Series) TableName() string {
	return "series"
}

type ListSeries struct {
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
type ListSeriesArr []ListSeries

func FormatSeries(Series Series) (list ListSeries) {
	list.ID = Series.ID
	list.RussianName = Series.RussianName
	list.Image = Series.Image
	list.Country = Series.Country
	list.YearFirst = Series.YearFirst
	list.YearLast = Series.YearLast
	list.AgeLimit = Series.AgeLimit
	list.Rating = Series.Rating
	return list
}

func (lists *ListSeriesArr) Convert(serials SeriesArr) ListSeriesArr {
	var lf ListSeriesArr
	for _, s := range serials {
		lf = append(lf, FormatSeries(s))
	}
	return lf
}
