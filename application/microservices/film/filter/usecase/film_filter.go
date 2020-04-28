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
	return f.filmRepo.FilterFilmData()
}

func (f *FilmFilter) FilterFilmList(fields map[string][]string) (*models.Films, bool) {
	return f.filmRepo.FilterFilmsList(fields)
}
