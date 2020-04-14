package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/jinzhu/gorm"
	"strconv"
)

var SeriesPerPage = 13

type PostgresForSerials struct {
	DB *gorm.DB
}

func NewPostgresForSeries(db *gorm.DB) series.Repository {
	return &PostgresForSerials{DB: db}
}

func (p PostgresForSerials) GetSeriesGenres(fid uint) (models.Genres, bool) {
	genres := &models.Genres{}
	db := p.DB.Table("kinopoisk.genres").Select("genres.name,genres.reference").
		Joins("join kinopoisk.series_genres on series_genres.genre_id=genres.id").
		Where("series_genres.series_id=?", fid).Order("genres.name").Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	return *genres, true
}

func (p PostgresForSerials) FilterSeriesList(fields map[string][]string) (*models.SeriesArr, bool) {
	series := &models.SeriesArr{}
	var db *gorm.DB
	var offset int
	var err error
	err = nil
	db = p.DB.Table("kinopoisk.series").
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
		offset = (offset - 1) * SeriesPerPage
	}
	if !pageOk || (err != nil) {
		return &models.SeriesArr{}, false
	}
	db = db.Group("kinopoisk.series.id").Order(order[0]).Offset(offset).Limit(SeriesPerPage).Find(series)

	err = db.Error
	if err != nil {
		fmt.Print(err, "\n")
		return &models.SeriesArr{}, false
	}
	return series, true
}

func (p PostgresForSerials) FilterSeriesData() (map[string]interface{}, bool) {
	genres := &models.Genres{}

	db := p.DB.Table("kinopoisk.genres").Order("genres.name").Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	var max, min int
	row := p.DB.Table("kinopoisk.series").Select("MAX(yearlast),MIN(yearfirst)").Row()
	err = row.Scan(&max, &min)
	if err != nil {
		return nil, false
	}
	resp := make(map[string]interface{})
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

func (p PostgresForSerials) GetSeriesByID(id uint) (models.Series, bool) {
	series := &models.Series{}
	db := p.DB.Table("kinopoisk.series").Find(series, id)
	err := db.Error
	if err != nil {
		return models.Series{}, false
	}
	return *series, true
}

func (p PostgresForSerials) GetSeriesSeasons(id uint) (models.Seasons, bool) {
	seasons := &models.Seasons{}
	db := p.DB.Table("kinopoisk.seasons").Where("seriesid = ?", id).Find(seasons)
	err := db.Error
	if err != nil {
		return models.Seasons{}, false
	}
	return *seasons, true
}

func (p PostgresForSerials) GetSeasonEpisodes(id uint) (models.Episodes, bool) {
	episodes := &models.Episodes{}
	db := p.DB.Table("kinopoisk.episodes").Where("seasonid = ?", id).Find(episodes)
	err := db.Error
	if err != nil {
		return models.Episodes{}, false
	}
	return *episodes, true
}
