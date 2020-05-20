package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/playlist"
	"github.com/go-park-mail-ru/2020_1_k-on/application/subscription"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"go.uber.org/zap"
)

type Subscription struct {
	plistRep playlist.Repository
	subsRep  subscription.Repository
	logger   *zap.Logger
}

func NewSubscription(plistRep playlist.Repository,
	subsRep subscription.Repository,
	logger *zap.Logger) subscription.UseCase {
	return &Subscription{
		plistRep: plistRep,
		subsRep:  subsRep,
		logger:   logger,
	}
}

func (s *Subscription) Subscribe(pid, userId uint) error {
	if pid == 0 || userId == 0 {
		return errors.NewInvalidArgument("wrong id")
	}

	play, err := s.plistRep.Get(pid)
	if err != nil {
		return err
	}
	if !play.Public {
		return errors.NewForbiddenError("playlist is private")
	}

	return s.subsRep.Subscribe(pid, userId)
}

func (s *Subscription) Unsubscribe(pid, userId uint) error {
	if pid == 0 || userId == 0 {
		return errors.NewInvalidArgument("wrong id")
	}

	return s.subsRep.Unsubscribe(pid, userId)
}

func (s *Subscription) Subscriptions(userId uint) (models.Playlists, error) {
	var plist models.Playlists

	if userId == 0 {
		return nil, errors.NewInvalidArgument("wrong id")
	}

	pidList, err := s.subsRep.Subscriptions(userId)
	if err != nil {
		return nil, err
	}

	for _, pid := range pidList {
		play, err := s.plistRep.Get(pid)
		if err != nil {
			return nil, err
		}

		if !play.Public && play.UserId != userId {
			continue
		}

		plist = append(plist, *play)
	}

	return plist, nil
}

func (s *Subscription) GetMainPagePlaylists(userId uint) (models.Playlists, error) {
	plist, err := s.subsRep.GetMainPagePlaylists(userId)
	if err != nil {
		return nil, err
	}

	for i, play := range plist {
		films, err := s.plistRep.GetFilms(play.Id)
		if err != nil {
			return nil, err
		}

		plist[i].Films = films

		series, err := s.plistRep.GetSeries(play.Id)
		if err != nil {
			return nil, err
		}

		plist[i].Series = series
	}

	return plist, nil
}
