package film

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
)

//Человеко читаемые методы, которые и будут вызываться в хендлерах в деливери

type Usecase interface {
	GetFilmsList(begin, end uint) (models.Films, bool)
	GetFilm(id uint) (models.Film, bool)
	CreateFilm(f models.Film) (models.Film, bool)
	GetFilmGenres(fid uint) (models.Genres, bool)
	Search(word string, query map[string][]string) (models.Films, bool)
	GetSimilarFilms(fid uint) (models.Films, bool)
	GetSimilarSeries(fid uint) (models.SeriesArr, bool)
}
