package subscription

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type UseCase interface {
	Subscribe(pid, userId uint) error
	Unsubscribe(pid, userId uint) error
	Subscriptions(userId uint) (models.Playlists, error)
}
