package client

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type IFilmFilterClient interface {
	GetFilteredFilms(fields map[string][]string) (models.Films, bool)
	GetFilterFields() (map[string]models.Genres, bool)
}
