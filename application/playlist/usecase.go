package playlist

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type UseCase interface {
	Create(playlist *models.Playlist) (*models.Playlist, error)
	AddFilm(pid, filmId, userId uint) error
	AddSeries(pid, seriesId, userId uint) error
	Get(pid uint, userId uint) (*models.Playlist, error)
	GetUserPlaylists(userId uint) (models.Playlists, error)
	GetUserPublicPlaylists(userId uint) (models.Playlists, error)
	Delete(pid uint, userId uint) error
	DeleteFilm(pid, filmId uint, userId uint) error
	DeleteSeries(pid, seriesId uint, userId uint) error
}
