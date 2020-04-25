package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/filter"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
)

type FilmFilter struct {
	filmRepo filter.Repository
}

func NewFilmFilter(filmRepo filter.Repository) filter.UseCase {
	return &FilmFilter{filmRepo: filmRepo}
}

func (f *FilmFilter) FilterFilmData() (map[string]models.Genres, bool) {
	data, ok := f.filmRepo.FilterFilmData()
	if !ok {
		return nil, false
	}
	return data, true
}

func (f *FilmFilter) FilterFilmList(fields map[string][]string) (models.Films, bool) {
	films, ok := f.filmRepo.FilterFilmsList(fields)
	if !ok {
		return models.Films{}, false
	}
	return *films, true
}
