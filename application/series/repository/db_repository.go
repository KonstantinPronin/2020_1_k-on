package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/jinzhu/gorm"
	"strings"
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

func (p PostgresForSerials) Search(word string, begin, end int) (models.SeriesArr, bool) {
	var seriesArr models.SeriesArr
	var query string
	words := strings.Split(word, " ")

	for _, str := range words {
		if query == "" {
			query = str
			continue
		}
		query = fmt.Sprintf("%s | %s", query, str)
	}

	err := p.DB.Table("kinopoisk.series").
		Where("textsearchable_index_col @@ to_tsquery('russian', ?)", query).
		Offset(begin).Limit(end).Find(&seriesArr).Error
	if err != nil {
		return nil, false
	}

	return seriesArr, true
}
