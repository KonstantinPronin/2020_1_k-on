package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/filter"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
)

type SeriesFilter struct {
	seriesRepo filter.Repository
}

func NewSeriesFilter(seriesRepo filter.Repository) filter.UseCase {
	return &SeriesFilter{seriesRepo: seriesRepo}
}

func (s *SeriesFilter) FilterSeriesList(fields map[string][]string) (*models.SeriesArr, bool) {
	return s.seriesRepo.FilterSeriesList(fields)
}

func (s *SeriesFilter) FilterSeriesData() (map[string]models.Genres, bool) {
	return s.seriesRepo.FilterSeriesData()
}
