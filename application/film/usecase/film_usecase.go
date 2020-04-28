package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/film"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
)

type filmUsecase struct {
	filmRepo film.Repository
}

func NewFilmUsecase(filmRepo film.Repository) film.Usecase {
	return &filmUsecase{filmRepo: filmRepo}
}

func (FU filmUsecase) GetFilmGenres(fid uint) (models.Genres, bool) {
	g, ok := FU.filmRepo.GetFilmGenres(fid)
	if !ok {
		return nil, false
	}
	return g, ok
}

func (FU filmUsecase) GetFilmsList(begin, end uint) (models.Films, bool) {
	films, ok := FU.filmRepo.GetFilmsArr(begin, end)
	if !ok {
		return models.Films{}, false
	}
	return *films, true
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

//
//func (FU filmUsecase) UploadImageFilm(id uint) models.Film {
//	f, _ := FU.filmRepo.GetById(id)
//	return *f
//}
//
//func (FU filmUsecase) GetImageFilm(id uint) string {
//	f, _ := FU.filmRepo.GetById(id)
//	return f.Image
//}
