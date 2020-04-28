package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/jinzhu/gorm"
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
		Joins("join kinopoisk.series_genres on series_genres.genre_ref=genres.reference").
		Where("series_genres.series_id=?", fid).Order("genres.name").Find(genres)
	err := db.Error
	if err != nil {
		return nil, false
	}
	return *genres, true
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
