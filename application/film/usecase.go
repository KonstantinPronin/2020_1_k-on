package film

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
)

//Человеко читаемые методы, которые и будут вызываться в хендлерах в деливери

type Usecase interface {
	GetFilmsList() (models.Films, bool)
	GetFilm(id uint) (models.Film, bool)
	CreateFilm(f models.Film) (models.Film, bool)
	FilterFilmList(fields map[string][]string) (models.Films, bool)
	FilterFilmData() (interface{}, bool)
	UploadImageFilm(id uint) models.Film //?? как правильно вести себя с картинками?
	GetImageFilm(id uint) string         //(model.File)????????????
}
