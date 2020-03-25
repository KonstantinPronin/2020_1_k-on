package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
)

type serialUsecase struct {
	serialRepo series.Repository
}

func NewSeriesUsecase(serialRepo series.Repository) series.Usecase {
	return &serialUsecase{serialRepo: serialRepo}
}

func (SU serialUsecase) GetSeriesSeasons(id uint) (models.Seasons, bool) {
	seasons, ok := SU.serialRepo.GetSeriesSeasons(id)
	if !ok {
		return models.Seasons{}, false
	}
	return seasons, true
}

func (SU serialUsecase) GetSeriesByID(id uint) (models.Series, bool) {
	serial, ok := SU.serialRepo.GetSeriesByID(id)
	if !ok {
		return models.Series{}, false
	}
	return serial, true
}

func (SU serialUsecase) GetSeasonEpisodes(id uint) (models.Episodes, bool) {
	series, ok := SU.serialRepo.GetSeasonEpisodes(id)
	if !ok {
		return models.Episodes{}, false
	}
	return series, true
}
