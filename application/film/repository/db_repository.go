package repository

import (
	_ "context"
	_ "errors"
	"github.com/go-park-mail-ru/2020_1_k-on/application/film"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/jinzhu/gorm"
)

//Интерфейсы запросов к бд

type PostgresForFilms struct {
	DB *gorm.DB
}

func NewPostgresForFilms(db *gorm.DB) film.Repository {
	return &PostgresForFilms{DB: db}
}

func (p PostgresForFilms) FilterFilmData() (map[string]interface{}, bool) {
	genres := &models.Genres{}

	db := p.DB.Table("kinopoisk.genres").Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	var max, min int
	row := db.Table("kinopoisk.films").Select("MAX(year),MIN(year)").Row()
	row.Scan(&max, &min)
	err = db.Error
	if err != nil {
		return nil, false
	}
	resp := make(map[string]interface{})
	resp["genres"] = genres
	resp["minyear"] = min
	resp["maxyear"] = max

	return resp, true
}

func (p PostgresForFilms) FilterFilmsList(fields map[string][]string) (*models.Films, bool) {
	films := &models.Films{}
	query := make(map[string]interface{})
	for key, val := range fields {
		query[key] = val[0]
	}
	db := p.DB.Table("kinopoisk.films").Where(query).Find(films)
	err := db.Error
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
