package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/serial"
)

type serialUsecase struct {
	serialRepo serial.Repository
}

func NewSerialUsecase(serialRepo serial.Repository) serial.Usecase {
	return &serialUsecase{serialRepo: serialRepo}
}

func (SU serialUsecase) GetSerialSeasons(id uint) (models.Seasons, bool) {
	seasons, ok := SU.serialRepo.GetSerialSeasons(id)
	if !ok {
		return models.Seasons{}, false
	}
	return seasons, true
}

func (SU serialUsecase) GetSerialByID(id uint) (models.Serial, bool) {
	serial, ok := SU.serialRepo.GetSerialByID(id)
	if !ok {
		return models.Serial{}, false
	}
	return serial, true
}

func (SU serialUsecase) GetSeasonSeries(id uint) (models.SeriesArr, bool) {
	series, ok := SU.serialRepo.GetSeasonSeries(id)
	if !ok {
		return models.SeriesArr{}, false
	}
	return series, true
}
