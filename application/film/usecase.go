package film

import (
	"2020_1_k-on/application/models"
)

//Человеко читаемые методы, которые и будут вызываться в хендлерах в деливери

type Usecase interface {
	GetFilmsList() models.Films
	GetFilm(id uint) (models.Film, bool)
	CreateFilm(f models.Film) models.Film
	UploadImageFilm(id uint) models.Film //?? как правильно вести себя с картинками?
	GetImageFilm(id uint) string         //(model.File)????????????
}
