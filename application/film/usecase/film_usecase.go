package usecase

import (
	"2020_1_k-on/application/film"
	"2020_1_k-on/application/models"
	"fmt"
)

type filmUsecase struct {
	filmRepo film.Repository
}

func NewFilmUsecase(filmRepo film.Repository) film.Usecase {
	return &filmUsecase{filmRepo: filmRepo}
}

func (FU filmUsecase) GetFilmsList() models.Films {
	b, e := uint(10), uint(0)
	films, ok := FU.filmRepo.GetFilmsArr(b, e)
	if !ok {
		fmt.Print(films)
	}
	return *films
}

func (FU filmUsecase) GetFilm(id uint) (models.Film, bool) {
	f, ok := FU.filmRepo.GetById(id)
	if !ok {
		return models.Film{}, false
	}
	return *f, true
}

func (FU filmUsecase) CreateFilm(f models.Film) (models.Film, bool) {
	var ok bool
	f, ok = FU.filmRepo.Create(&f)
	return f, ok
}

func (FU filmUsecase) UploadImageFilm(id uint) models.Film {
	f, _ := FU.filmRepo.GetById(id)
	return *f
}

func (FU filmUsecase) GetImageFilm(id uint) string {
	f, _ := FU.filmRepo.GetById(id)
	return f.Image
}
