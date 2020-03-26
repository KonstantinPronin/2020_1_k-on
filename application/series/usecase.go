package series

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
)

//Человеко читаемые методы, которые и будут вызываться в хендлерах в деливери

type Usecase interface {
	GetSeriesByID(id uint) (models.Series, bool)
	GetSeriesSeasons(id uint) (models.Seasons, bool)
	GetSeasonEpisodes(id uint) (models.Episodes, bool)
}
