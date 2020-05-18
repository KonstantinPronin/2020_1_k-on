package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/film"
	"github.com/go-park-mail-ru/2020_1_k-on/application/film/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"strconv"
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

func (FU filmUsecase) Search(word string, query map[string][]string) (models.Films, bool) {
	begin := 0

	page, ok := query["page"]
	if ok {
		var err error
		begin, err = strconv.Atoi(page[0])
		if err == nil {
			begin = (begin - 1) * repository.FilmPerPage
		}
		if begin < 0 {
			begin = 0
		}
	}

	return FU.filmRepo.Search(word, begin, repository.FilmPerPage)
}

func (FU filmUsecase) GetSimilarFilms(fid uint) (models.Films, bool) {
	f, ok := FU.filmRepo.GetSimilarFilms(fid)
	if !ok {
		return models.Films{}, false
	}
	return f, true
}

func (FU filmUsecase) GetSimilarSeries(fid uint) (models.SeriesArr, bool) {
	ser, ok := FU.filmRepo.GetSimilarSeries(fid)
	if !ok {
		return models.SeriesArr{}, false
	}
	return ser, true
}
