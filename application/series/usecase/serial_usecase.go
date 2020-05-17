package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series/repository"
	"strconv"
)

type serialUsecase struct {
	serialRepo series.Repository
}

func NewSeriesUsecase(serialRepo series.Repository) series.Usecase {
	return &serialUsecase{serialRepo: serialRepo}
}

func (SU serialUsecase) GetSeriesGenres(fid uint) (models.Genres, bool) {
	g, ok := SU.serialRepo.GetSeriesGenres(fid)
	if !ok {
		return nil, false
	}
	return g, ok
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

func (SU serialUsecase) Search(word string, query map[string][]string) (models.SeriesArr, bool) {
	begin := 0
	page, ok := query["page"]
	if ok {
		var err error
		begin, err = strconv.Atoi(page[0])
		if err == nil {
			begin = (begin - 1) * repository.SeriesPerPage
		}
		if begin < 0 {
			begin = 0
		}
	}

	return SU.serialRepo.Search(word, begin, repository.SeriesPerPage)
}
