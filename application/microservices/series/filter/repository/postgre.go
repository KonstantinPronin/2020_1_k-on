package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/filter"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series/repository"
	"github.com/jinzhu/gorm"
	"strconv"
)

type SeriesFiltersDb struct {
	DB *gorm.DB
}

func NewSeriesFiltersDb(DB *gorm.DB) filter.Repository {
	return &SeriesFiltersDb{DB: DB}
}

func (s *SeriesFiltersDb) FilterSeriesList(fields map[string][]string) (*models.SeriesArr, bool) {
	series := &models.SeriesArr{}
	var db *gorm.DB
	var offset int
	var err error
	err = nil
	db = s.DB.Table("kinopoisk.series").
		Joins("join kinopoisk.series_genres on kinopoisk.series_genres.series_id=kinopoisk.series.id")

	for key, val := range fields {
		if val[0] == "ALL" {
			delete(fields, key)
		} else {
			if key == "year" {
				db = db.Where("yearfirst <= ? AND yearlast >= ?", val[0], val[0])
			}
			if key == "genre" {
				db = db.Where("series_genres.genre_ref = ? ", val[0])
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
		offset = (offset - 1) * repository.SeriesPerPage
	}
	if !pageOk || (err != nil) {
		return &models.SeriesArr{}, false
	}
	db = db.Group("kinopoisk.series.id").Order(order[0]).Offset(offset).Limit(repository.SeriesPerPage).Find(series)

	err = db.Error
	if err != nil {
		fmt.Print(err, "\n")
		return &models.SeriesArr{}, false
	}
	return series, true
}

func (s *SeriesFiltersDb) FilterSeriesData() (map[string]models.Genres, bool) {
	genres := &models.Genres{}

	db := s.DB.Table("kinopoisk.genres").Order("genres.name").Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	var max, min int
	row := s.DB.Table("kinopoisk.series").Select("MAX(yearlast),MIN(yearfirst)").Row()
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
		models.Genre{"По рейтингу", "rating"},
		models.Genre{"По рейтингу IMDb", "imdbrating"},
	}
	resp["year"] = models.Genres{
		models.Genre{"Все годы", "%"},
		models.Genre{"maxyear", strconv.Itoa(max)},
		models.Genre{"minyear", strconv.Itoa(min)},
	}
	return resp, true
}
