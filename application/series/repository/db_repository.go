package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/series"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/util"
	"github.com/jinzhu/gorm"
)

var SeriesPerPage = 13
var SimilarCount = 10

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
	query := util.PlainTextToQuery(word)
	word = fmt.Sprintf("%%%s%%", word)

	err := p.DB.Table("kinopoisk.series").
		Where("textsearchable_index_col @@ to_tsquery('russian', ?)  or russianname ilike ?", query, word).
		Offset(begin).Limit(end).Find(&seriesArr).Error
	if err != nil {
		return nil, false
	}

	return seriesArr, true
}

func (p PostgresForSerials) GetSimilarFilms(sid uint) (models.Films, bool) {
	films := &models.Films{}
	var db *gorm.DB
	var err error
	err = nil
	db = p.DB.Table("kinopoisk.films f2").Select("f2.*").
		Joins("join (?) sub on f2.russianname = sub.russianname",
			p.DB.Table("kinopoisk.series_playlist sp1").
				Select("f1.russianname, count(fp2.film_id)").
				Joins("join kinopoisk.film_playlist fp2 on sp1.playlist_id = fp2.playlist_id").
				Joins("join kinopoisk.films f1 on fp2.film_id = f1.id").
				Where("sp1.series_id = ?", sid).Group("f1.russianname").
				SubQuery()).
		Order("sub.count desc").Limit(SimilarCount).Find(films)
	err = db.Error
	if err != nil {
		fmt.Print(err, "\n")
		return models.Films{}, false
	}
	return *films, true
}

func (p PostgresForSerials) GetSimilarSeries(sid uint) (models.SeriesArr, bool) {
	series := &models.SeriesArr{}
	var db *gorm.DB
	var err error
	err = nil
	db = p.DB.Table("kinopoisk.series s2").Select("s2.*").
		Joins("join (?) sub on s2.russianname = sub.russianname",
			p.DB.Table("kinopoisk.series_playlist sp1").
				Select("s1.russianname, count(sp2.series_id)").
				Joins("join kinopoisk.series_playlist sp2 on sp1.playlist_id = sp2.playlist_id").
				Joins("join kinopoisk.series s1 on sp2.series_id = s1.id").
				Where("sp1.series_id = ?", sid).Group("s1.russianname").
				SubQuery()).
		Where("s2.id <> ?", sid).
		Order("sub.count desc").Limit(SimilarCount).Find(series)
	err = db.Error
	if err != nil {
		fmt.Print(err, "\n")
		return models.SeriesArr{}, false
	}
	return *series, true
}
