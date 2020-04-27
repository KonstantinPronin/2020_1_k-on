package playlist

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type Repository interface {
	Create(playlist *models.Playlist) (*models.Playlist, error)
	AddFilm(pid, filmId uint) error
	AddSeries(pid, seriesId uint) error

	Get(pid uint) (*models.Playlist, error)
	GetUserPlaylists(userId uint) (models.Playlists, error)
	GetUserPublicPlaylists(userId uint) (models.Playlists, error)
	GetPlaylistsWithoutSer(sid, userId uint) (models.Playlists, error)
	GetPlaylistsWithoutFilm(fid, userId uint) (models.Playlists, error)
	GetAdminPlaylists() (models.Playlists, error)

	Delete(pid uint) error
	DeleteFilm(pid, filmId uint) error
	DeleteSeries(pid, seriesId uint) error
	IsCreator(pid, userId uint) (bool, error)
}
