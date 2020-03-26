package models

type Review struct {
	Id     uint   `json:"id" gorm:"primary_key"`
	Rating int    `json:"rating" gorm:"column:rating"`
	Body   string `json:"body" gorm:"column:body"`
	UserId uint   `json:"userId" gorm:"column:userid"`
	FilmId uint   `json:"filmId" gorm:"column:filmid"`
}
