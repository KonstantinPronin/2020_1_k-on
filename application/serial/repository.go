package serial

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type Repository interface {
	GetSerialByID(id uint) (models.Serial, bool)
	GetSerialSeasons(id uint) (models.Seasons, bool)
	GetSeasonSeries(id uint) (models.SeriesArr, bool)
}
