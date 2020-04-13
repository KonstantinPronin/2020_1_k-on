package repository

import (
	_ "context"
	_ "errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_1_k-on/application/film"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/jinzhu/gorm"
	"strconv"
)

//Интерфейсы запросов к бд

var FilmPerPage = 13

type PostgresForFilms struct {
	DB *gorm.DB
}

func NewPostgresForFilms(db *gorm.DB) film.Repository {
	return &PostgresForFilms{DB: db}
}

func (p PostgresForFilms) GetFilmGenres(fid uint) (models.Genres, bool) {
	genres := &models.Genres{}
	db := p.DB.Table("kinopoisk.genres").Select("genres.name,genres.reference").
		Joins("join kinopoisk.films_genres on films_genres.genre_ref=genres.reference").
		Where("films_genres.film_id=?", fid).Order("genres.name").Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	return *genres, true
}

func (p PostgresForFilms) FilterFilmData() (map[string]interface{}, bool) {
	genres := &models.Genres{}

	db := p.DB.Table("kinopoisk.genres").Order("name").Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	var max, min int
	row := p.DB.Table("kinopoisk.films").Select("MAX(films.year),MIN(films.year)").Row()
	err = row.Scan(&max, &min)
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
	db = p.DB.Table("kinopoisk.films").
		Joins("join kinopoisk.films_genres on kinopoisk.films_genres.film_id=kinopoisk.films.id")

	for key, val := range fields {
		if val[0] == "ALL" {
			delete(fields, key)
		} else {
			if key == "year" {
				db = db.Where("year = ?", val[0])
			}
			if key == "genre" {
				db = db.Where("films_genres.genre_ref = ? ", val[0])
			}
		}
	}

	order, ok := fields["order"]
	if ok {
		order[0] = order[0] + " DESC"
		delete(fields, "order")
	} else {
		order = []string{"-rating"}
	}

	page, pageOk := fields["page"]
	if pageOk {
		delete(fields, "page")
		offset, err = strconv.Atoi(page[0])
		offset = (offset - 1) * FilmPerPage
	}
	if !pageOk || (err != nil) {
		return &models.Films{}, false
	}
	db = db.Group("kinopoisk.films.id").Order(order[0]).Offset(offset).Limit(FilmPerPage).Find(films)

	err = db.Error
	if err != nil {
		fmt.Print(err, "\n")
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
