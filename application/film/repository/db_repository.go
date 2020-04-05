package repository

import (
	_ "context"
	_ "errors"
	"github.com/go-park-mail-ru/2020_1_k-on/application/film"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/jinzhu/gorm"
	"strconv"
)

//Интерфейсы запросов к бд

var FilmPerPage = 10

type PostgresForFilms struct {
	DB *gorm.DB
}

func NewPostgresForFilms(db *gorm.DB) film.Repository {
	return &PostgresForFilms{DB: db}
}

func (p PostgresForFilms) GetFilmGenres(fid uint) (models.Genres, bool) {
	genres := &models.Genres{}
	db := p.DB.Table("kinopoisk.genres").Select("genres.id,genres.name,genres.reference").
		Joins("join kinopoisk.films_genres on films_genres.genre_id=genres.id").Where("films_genres.film_id=?", fid).Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	//db.Close()
	return *genres, true
}

func (p PostgresForFilms) FilterFilmData() (map[string]interface{}, bool) {
	genres := &models.Genres{}

	db := p.DB.Table("kinopoisk.genres").Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	//db.Close()
	var max, min int
	row := db.Table("kinopoisk.films").Select("MAX(year),MIN(year)").Row()
	row.Scan(&max, &min)
	err = db.Error
	//db.Close()
	if err != nil {
		return nil, false
	}
	resp := make(map[string]interface{})
	filters := make(map[string]interface{})
	filters["minyear"] = min
	filters["maxyear"] = max
	resp["genres"] = genres
	resp["filters"] = filters

	return resp, true
}

func (p PostgresForFilms) FilterFilmsList(fields map[string][]string) (*models.Films, bool) {
	films := &models.Films{}
	var db *gorm.DB
	var offset int
	var err error
	err = nil
	query := make(map[string]interface{})
	order, ok := fields["order"]
	if ok {
		delete(fields, "order")
	}
	page, pok := fields["page"]
	if pok {
		delete(fields, "page")
		offset, err = strconv.Atoi(page[0])
		offset = (offset - 1) * FilmPerPage
	}
	if !pok || (err != nil) {
		return &models.Films{}, false
	}
	for key, val := range fields {
		query[key] = val[0]
	}
	if ok {
		db = p.DB.Table("kinopoisk.films").Where(query).Order(order[0]).Offset(offset).Limit(FilmPerPage).Find(films)
	} else {
		db = p.DB.Table("kinopoisk.films").Where(query).Offset(offset).Limit(FilmPerPage).Find(films)
	}
	err = db.Error
	if err != nil {
		return &models.Films{}, false
	}
	return films, true
}

func (p PostgresForFilms) GetFilmsArr(begin, end uint) (*models.Films, bool) {
	films := &models.Films{}
	db := p.DB.Table("kinopoisk.films").Offset(end).Limit(begin).Find(films)
	err := db.Error
	if err != nil {
		return &models.Films{}, false
	}
	return films, true
}

func (p PostgresForFilms) Create(film *models.Film) (models.Film, bool) {
	db := p.DB.Table("kinopoisk.films").Create(film)
	err := db.Error
	if err != nil {
		return models.Film{}, false
	}
	return *film, true
}

func (p PostgresForFilms) GetById(id uint) (*models.Film, bool) {
	f := &models.Film{}
	db := p.DB.Table("kinopoisk.films").Find(f, id)
	err := db.Error
	if err != nil {
		return &models.Film{}, false
	}
	return f, true
}

func (p PostgresForFilms) GetByName(name string) (*models.Film, bool) {
	f := &models.Film{}
	db := p.DB.Table("kinopoisk.films").Where("englishname = ?", name).First(&f)
	err := db.Error
	if err != nil {
		return &models.Film{}, false
	}
	return f, true
}
