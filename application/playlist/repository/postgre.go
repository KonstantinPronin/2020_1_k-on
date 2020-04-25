package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/playlist"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type PlaylistDatabase struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewPlaylistDatabase(db *gorm.DB, logger *zap.Logger) playlist.Repository {
	return &PlaylistDatabase{
		conn:   db,
		logger: logger,
	}
}

func (p *PlaylistDatabase) Create(playlist *models.Playlist) (*models.Playlist, error) {
	err := p.conn.Table("kinopoisk.playlists").Create(playlist).Error
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (p *PlaylistDatabase) AddFilm(pid, filmId uint) error {
	return p.conn.Table("kinopoisk.film_playlist").Create(&models.FilmToPlaylist{
		Pid:    pid,
		FilmId: filmId,
	}).Error
}

func (p *PlaylistDatabase) AddSeries(pid, seriesId uint) error {
	return p.conn.Table("kinopoisk.series_playlist").Create(&models.SeriesToPlaylist{
		Pid:      pid,
		SeriesId: seriesId,
	}).Error
}

func (p *PlaylistDatabase) Get(pid uint) (*models.Playlist, error) {
	play := new(models.Playlist)
	err := p.conn.Table("kinopoisk.playlists").
		Where("id = ?", pid).First(play).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewNotFoundError(err.Error())
		}
		return nil, err
	}

	play.Films, err = p.GetFilms(pid)
	if err != nil {
		return nil, err
	}

	play.Series, err = p.GetSeries(pid)
	if err != nil {
		return nil, err
	}

	return play, nil
}

func (p *PlaylistDatabase) GetUserPlaylists(userId uint) (models.Playlists, error) {
	var plist models.Playlists

	rows, err := p.conn.Table("kinopoisk.playlists").
		Select("id").Where("user_id = ?", userId).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var pid uint
		err := rows.Scan(&pid)
		if err != nil {
			return nil, err
		}

		play, err := p.Get(pid)
		if err != nil {
			return nil, err
		}

		plist = append(plist, *play)
	}

	return plist, nil
}

func (p *PlaylistDatabase) Delete(pid uint) error {
	return p.conn.Table("kinopoisk.playlists").
		Where("id = ?", pid).
		Delete(&models.Playlist{}).Error
}

func (p *PlaylistDatabase) DeleteFilm(pid, filmId uint) error {
	return p.conn.Table("kinopoisk.film_playlist").
		Where("playlist_id = ? and film_id = ?", pid, filmId).
		Delete(&models.FilmToPlaylist{}).Error
}

func (p *PlaylistDatabase) DeleteSeries(pid, seriesId uint) error {
	return p.conn.Table("kinopoisk.series_playlist").
		Where("playlist_id = ? and series_id = ?", pid, seriesId).
		Delete(&models.SeriesToPlaylist{}).Error
}

func (p *PlaylistDatabase) IsCreator(pid, userId uint) (bool, error) {
	play := new(models.Playlist)
	err := p.conn.Table("kinopoisk.playlists").
		Where("id = ? and user_id = ?", pid, userId).First(play).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (p *PlaylistDatabase) GetFilms(pid uint) (models.ListsFilm, error) {
	var films models.ListsFilm

	rows, err := p.conn.Table("kinopoisk.films f").
		Select("f.id, f.maingenre, f.russianname, f.image, f.country, f.year, f.agelimit, f.rating").
		Joins("inner join kinopoisk.film_playlist p on f.id = p.film_id").
		Where("p.playlist_id = ?", pid).Rows()
	if err != nil {
		return nil, err
	}

	film := new(models.ListFilm)
	for rows.Next() {
		err := rows.Scan(&film.ID, &film.MainGenre, &film.RussianName, &film.Image, &film.Country, &film.Year, &film.AgeLimit, &film.Rating)
		if err != nil {
			return nil, err
		}

		film.Type = "films"
		films = append(films, *film)
	}

	return films, err
}

func (p *PlaylistDatabase) GetSeries(pid uint) (models.ListSeriesArr, error) {
	var series models.ListSeriesArr

	rows, err := p.conn.Table("kinopoisk.series s").
		Select("s.id, s.maingenre, s.russianname, s.image, s.country, s.yearfirst, s.yearlast, s.agelimit, s.rating").
		Joins("inner join kinopoisk.series_playlist p on p.series_id = s.id").
		Where("p.playlist_id = ?", pid).Rows()
	if err != nil {
		return nil, err
	}

	s := new(models.ListSeries)
	for rows.Next() {
		err := rows.Scan(&s.ID, &s.MainGenre, &s.RussianName, &s.Image, &s.Country, &s.YearFirst, &s.YearLast, &s.AgeLimit, &s.Rating)
		if err != nil {
			return nil, err
		}
		s.Type = "series"
		series = append(series, *s)
	}

	return series, nil
}
