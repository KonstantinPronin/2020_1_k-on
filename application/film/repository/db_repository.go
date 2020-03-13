package repository

import (
	"2020_1_k-on/application/film"
	"2020_1_k-on/application/models"
	_ "context"
	_ "errors"
	"github.com/jinzhu/gorm"
)

//Интерфейсы запросов к бд

type PostgresForFilms struct {
	DB *gorm.DB
}

func NewPostgresForFilms(db *gorm.DB) film.Repository {
	return &PostgresForFilms{DB: db}
}

func (p PostgresForFilms) Create(film *models.Film) (models.Film, bool) {
	db := p.DB.Create(film)
	err := db.Error
	if err != nil {
		return models.Film{}, false
	}
	return *film, true
}

func (p PostgresForFilms) GetById(id uint) (*models.Film, bool) {
	film := &models.Film{}
	db := p.DB.Find(film, id)
	err := db.Error
	if err != nil {
		return &models.Film{}, false
	}
	return film, true
}

func (p PostgresForFilms) GetByName(name string) (*models.Film, bool) {
	film := &models.Film{}
	db := p.DB.Where("name = ?", name).First(&film)
	err := db.Error
	if err != nil {
		return &models.Film{}, false
	}
	return film, true
}

func (p PostgresForFilms) GetFilmsArr(begin, end uint) (*models.Films, bool) {
	films := &models.Films{}
	db := p.DB.Offset(end).Limit(begin).Find(films)
	err := db.Error
	if err != nil {
		return &models.Films{}, false
	}
	return films, true
}
