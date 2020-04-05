package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/jinzhu/gorm"
	"strconv"
)

var SeriesPerPage = 10

type PostgresForSerials struct {
	DB *gorm.DB
}

func NewPostgresForSeries(db *gorm.DB) series.Repository {
	return &PostgresForSerials{DB: db}
}

func (p PostgresForSerials) GetSeriesGenres(fid uint) (models.Genres, bool) {
	genres := &models.Genres{}
	db := p.DB.Table("kinopoisk.genres").Select("genres.id,genres.name,genres.reference").
		Joins("join kinopoisk.series_genres on series_genres.genre_id=genres.id").Where("series_genres.series_id=?", fid).Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	db.Close()
	return *genres, true
}

func (p PostgresForSerials) FilterSeriesList(fields map[string][]string) (*models.SeriesArr, bool) {
	series := &models.SeriesArr{}
	var db *gorm.DB
	var offset int
	var err error = nil
	query := make(map[string]interface{})
	order, ok := fields["order"]
	if ok {
		delete(fields, "order")
	}
	page, pok := fields["page"]
	if pok {
		delete(fields, "page")
		offset, err = strconv.Atoi(page[0])
		offset = (offset - 1) * SeriesPerPage
	}
	if !pok || (err != nil) {
		return &models.SeriesArr{}, false
	}
	for key, val := range fields {
		query[key] = val[0]
	}
	if ok {
		db = p.DB.Table("kinopoisk.series").Where(query).Order(order[0]).Offset(offset).Limit(SeriesPerPage).Find(series)
	} else {
		db = p.DB.Table("kinopoisk.series").Where(query).Offset(offset).Limit(SeriesPerPage).Find(series)
	}
	err = db.Error
	if err != nil {
		return &models.SeriesArr{}, false
	}
	return series, true
}

func (p PostgresForSerials) FilterSeriesData() (map[string]interface{}, bool) {
	genres := &models.Genres{}

	db := p.DB.Table("kinopoisk.genres").Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	var max, min int
	row := db.Table("kinopoisk.series").Select("MAX(yearlast),MIN(yearfirst)").Row()
	row.Scan(&max, &min)
	err = db.Error
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
