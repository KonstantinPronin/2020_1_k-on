package filter

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type UseCase interface {
	FilterFilmList(fields map[string][]string) (*models.Films, bool)
	FilterFilmData() (map[string]models.Genres, bool)
}
