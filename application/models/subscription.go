package models

type Subscription struct {
	Id     uint `gorm:"primary_key"`
	Pid    uint `gorm:"column:playlist_id"`
	UserId uint `gorm:"column:user_id"`
}
