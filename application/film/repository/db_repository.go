package repository

import (
	_ "context"
	_ "errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_1_k-on/application/film"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/util"
	"github.com/jinzhu/gorm"
)

//Интерфейсы запросов к бд

var FilmPerPage = 13
var SimilarCount = 10

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

func (p *PostgresForFilms) Search(word string, begin, end int) (models.Films, bool) {
	var films models.Films
	query := util.PlainTextToQuery(word)
	word = fmt.Sprintf("%%%s%%", word)

	err := p.DB.Table("kinopoisk.films").
		Where("textsearchable_index_col @@ to_tsquery('russian', ?) or russianname ilike ?", query, word).
		Offset(begin).Limit(end).Find(&films).Error
	if err != nil {
		return nil, false
	}

	return films, true
}

func (p PostgresForFilms) GetSimilarFilms(fid uint) (models.Films, bool) {
	films := &models.Films{}
	var db *gorm.DB
	var err error
	db = p.DB.Table("kinopoisk.films f2").Select("f2.*").
		Joins("join (?) sub on f2.russianname = sub.russianname",
			p.DB.Table("kinopoisk.film_playlist fp1").
				Select("f1.russianname, count(fp2.film_id)").
				Joins("join kinopoisk.film_playlist fp2 on fp1.playlist_id = fp2.playlist_id").
				Joins("join kinopoisk.films f1 on fp2.film_id = f1.id").
				Where("fp1.film_id = ?", fid).Group("f1.russianname").
				SubQuery()).
		Where("f2.id <> ?", fid).
		Order("sub.count desc").Limit(SimilarCount).Find(films)
	err = db.Error
	if err != nil {
		fmt.Print(err, "\n")
		return models.Films{}, false
	}
	return *films, true
}

func (p PostgresForFilms) GetSimilarSeries(fid uint) (models.SeriesArr, bool) {
	series := &models.SeriesArr{}
	var db *gorm.DB
	var err error
	db = p.DB.Table("kinopoisk.series s2").Select("s2.*").
		Joins("join (?) sub on s2.russianname = sub.russianname",
			p.DB.Table("kinopoisk.film_playlist fp1").
				Select("s1.russianname, count(sp2.series_id)").
				Joins("join kinopoisk.series_playlist sp2 on fp1.playlist_id = sp2.playlist_id").
				Joins("join kinopoisk.series s1 on sp2.series_id = s1.id").
				Where("fp1.film_id = ?", fid).Group("s1.russianname").
				SubQuery()).
		Order("sub.count desc").Limit(SimilarCount).Find(series)
	err = db.Error
	if err != nil {
		fmt.Print(err, "\n")
		return models.SeriesArr{}, false
	}
	return *series, true
}
