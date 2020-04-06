package models

//СЕРИАЛЫ!

type Series struct {
	ID              uint    `json:"id" gorm:"primary_key"`
	MainGenre       string  `json:"mainGenre" gorm:"column:maingenre"` //русское
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
	YearFirst       int     `json:"yearFirst" gorm:"column:yearfirst"` //начало
	YearLast        int     `json:"yearLast" gorm:"column:yearlast"`   //0 - "не закончился"
	AgeLimit        int     `json:"ageLimit,omitempty" gorm:"column:agelimit"`
}

//easyjson:json
type SeriesArr []Series

func (s *Series) TableName() string {
	return "series"
}

type ListSeries struct {
	ID          uint    `json:"id" gorm:"primary_key"`
	RussianName string  `json:"russianName" gorm:"column:russianname"`
	Image       string  `json:"image,omitempty"`
	Country     string  `json:"country"`
	YearFirst   int     `json:"yearFirst" gorm:"column:yearfirst"` //начало
	YearLast    int     `json:"yearLast" gorm:"column:yearlast"`   //0-  "не закончился"
	AgeLimit    int     `json:"ageLimit,omitempty" gorm:"column:agelimit"`
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
