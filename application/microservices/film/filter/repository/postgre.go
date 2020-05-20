package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_k-on/application/film/repository"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/filter"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/jinzhu/gorm"
	"strconv"
)

type FilmFiltersDb struct {
	DB *gorm.DB
}

func NewFilmFiltersDb(db *gorm.DB) filter.Repository {
	return &FilmFiltersDb{DB: db}
}

func (p *FilmFiltersDb) FilterFilmData() (map[string]models.Genres, bool) {
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
	resp := make(map[string]models.Genres)
	g := models.Genres{}
	g = append(g, models.Genre{
		Name:      "Все жанры",
		Reference: "%",
	})
	g = append(g, *genres...)
	resp["genre"] = g
	resp["order"] = models.Genres{
		models.Genre{Name: "По рейтингу", Reference: "rating"},
		models.Genre{Name: "По рейтингу IMDb", Reference: "imdbrating"},
	}
	resp["year"] = models.Genres{
		models.Genre{Name: "Все годы", Reference: "%"},
		models.Genre{Name: "maxyear", Reference: strconv.Itoa(max)},
		models.Genre{Name: "minyear", Reference: strconv.Itoa(min)},
	}

	return resp, true
}

func (p *FilmFiltersDb) FilterFilmsList(fields map[string][]string) (*models.Films, bool) {
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
		offset = (offset - 1) * repository.FilmPerPage
	}
	if !pageOk || (err != nil) {
		return &models.Films{}, false
	}
	db = db.Group("kinopoisk.films.id").Order(order[0]).Offset(offset).Limit(repository.FilmPerPage).Find(films)

	err = db.Error
	if err != nil {
		fmt.Print(err, "\n")
		return &models.Films{}, false
	}
	return films, true
}
