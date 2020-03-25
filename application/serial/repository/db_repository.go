package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/serial"
	"github.com/jinzhu/gorm"
)

type PostgresForSerials struct {
	DB *gorm.DB
}

func NewPostgresForserial(db *gorm.DB) serial.Repository {
	return &PostgresForSerials{DB: db}
}

func (p PostgresForSerials) GetSerialByID(id uint) (models.Serial, bool) {
	serial := &models.Serial{}
	db := p.DB.Table("kinopoisk.serials").Find(serial, id)
	err := db.Error
	if err != nil {
		return models.Serial{}, false
	}
	return *serial, true
}

func (p PostgresForSerials) GetSerialSeasons(id uint) (models.Seasons, bool) {
	seasons := &models.Seasons{}
	db := p.DB.Table("kinopoisk.seasons").Where("serialid = ?", id).Find(seasons)
	err := db.Error
	if err != nil {
		return models.Seasons{}, false
	}
	return *seasons, true
}

func (p PostgresForSerials) GetSeasonSeries(id uint) (models.SeriesArr, bool) {
	series := &models.SeriesArr{}
	db := p.DB.Table("kinopoisk.series").Where("seasonid = ?", id).Find(series)
	err := db.Error
	if err != nil {
		return models.SeriesArr{}, false
	}
	return *series, true
}
