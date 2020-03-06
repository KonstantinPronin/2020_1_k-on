package usecase

import (
	"2020_1_k-on/application/film"
	"2020_1_k-on/application/film/repository"
	"2020_1_k-on/application/models"
)

type filmUsecase struct {
	filmRepo film.Repository
}

func NewUserUsecase(filmRepo *repository.FilmsList) film.Usecase {
	return &filmUsecase{filmRepo: filmRepo}
}

func (FU filmUsecase) GetFilmsList() models.Films {
	var fms models.Films
	counter := 1
	f, ok := FU.filmRepo.GetById(uint(counter))
	for ok {
		fms = append(fms, *f)
		counter += 1
		f, ok = FU.filmRepo.GetById(uint(counter))
	}
	return fms
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
