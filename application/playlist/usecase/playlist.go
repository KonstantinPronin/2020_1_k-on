package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/playlist"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"go.uber.org/zap"
)

type Playlist struct {
	rep    playlist.Repository
	logger *zap.Logger
}

func NewPlaylist(rep playlist.Repository, logger *zap.Logger) playlist.UseCase {
	return &Playlist{
		rep:    rep,
		logger: logger,
	}
}

func (p *Playlist) Create(playlist *models.Playlist) (*models.Playlist, error) {
	if playlist.Name == "" || playlist.UserId == 0 {
		return nil, errors.NewInvalidArgument("wrong playlist parameters")
	}

	return p.rep.Create(playlist)
}

func (p *Playlist) AddFilm(pid, filmId, userId uint) error {
	if pid == 0 || filmId == 0 {
		return errors.NewInvalidArgument("wrong id")
	}

	ok, err := p.rep.IsCreator(pid, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.NewForbiddenError("user is not playlist creator")
	}

	return p.rep.AddFilm(pid, filmId)
}

func (p *Playlist) AddSeries(pid, seriesId, userId uint) error {
	if pid == 0 || seriesId == 0 {
		return errors.NewInvalidArgument("wrong id")
	}

	ok, err := p.rep.IsCreator(pid, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.NewForbiddenError("user is not playlist creator")
	}

	return p.rep.AddSeries(pid, seriesId)
}

func (p *Playlist) Get(pid, userId uint) (*models.Playlist, error) {
	if pid == 0 {
		return nil, errors.NewInvalidArgument("wrong id")
	}

	play, err := p.rep.Get(pid)
	if err != nil {
		return nil, err
	}

	if play.Public == true || play.UserId == userId {
		return play, nil
	} else {
		return nil, errors.NewForbiddenError("playlist is not public")
	}
}

func (p *Playlist) GetUserPlaylists(userId uint) (models.Playlists, error) {
	if userId == 0 {
		return nil, errors.NewInvalidArgument("wrong id")
	}

	return p.rep.GetUserPlaylists(userId)
}

func (p *Playlist) GetUserPublicPlaylists(userId uint) (models.Playlists, error) {
	if userId == 0 {
		return nil, errors.NewInvalidArgument("wrong id")
	}

	return p.rep.GetUserPublicPlaylists(userId)
}

func (p *Playlist) Delete(pid, userId uint) error {
	if pid == 0 {
		return errors.NewInvalidArgument("wrong id")
	}

	ok, err := p.rep.IsCreator(pid, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.NewForbiddenError("user is not playlist creator")
	}

	return p.rep.Delete(pid)
}

func (p *Playlist) DeleteFilm(pid, filmId, userId uint) error {
	if pid == 0 || filmId == 0 {
		return errors.NewInvalidArgument("wrong id")
	}

	ok, err := p.rep.IsCreator(pid, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.NewForbiddenError("user is not playlist creator")
	}

	return p.rep.DeleteFilm(pid, filmId)
}

func (p *Playlist) DeleteSeries(pid, seriesId, userId uint) error {
	if pid == 0 || seriesId == 0 {
		return errors.NewInvalidArgument("wrong id")
	}

	ok, err := p.rep.IsCreator(pid, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.NewForbiddenError("user is not playlist creator")
	}

	return p.rep.DeleteSeries(pid, seriesId)
}
