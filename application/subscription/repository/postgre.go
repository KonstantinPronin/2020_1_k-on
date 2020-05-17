package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/subscription"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"strings"
)

type SubscriptionDatabase struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewSubscriptionDatabase(db *gorm.DB, logger *zap.Logger) subscription.Repository {
	return &SubscriptionDatabase{
		conn:   db,
		logger: logger,
	}
}

func (s *SubscriptionDatabase) Subscribe(pid, userId uint) error {
	err := s.conn.Table("kinopoisk.subscriptions").Create(&models.Subscription{
		Pid:    pid,
		UserId: userId,
	}).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.NewInvalidArgument("user already subscribed to playlist")
		}
		return err
	}

	return nil
}

func (s *SubscriptionDatabase) Unsubscribe(pid, userId uint) error {
	return s.conn.Table("kinopoisk.subscriptions").
		Where("playlist_id = ? and user_id = ?", pid, userId).
		Delete(&models.Subscription{}).Error
}

func (s *SubscriptionDatabase) Subscriptions(userId uint) ([]uint, error) {
	var pidList []uint

	rows, err := s.conn.Table("kinopoisk.subscriptions").
		Select("playlist_id").
		Where("user_id = ?", userId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pid uint
		err := rows.Scan(&pid)
		if err != nil {
			return nil, err
		}

		pidList = append(pidList, pid)
	}

	return pidList, nil
}

func (s *SubscriptionDatabase) GetMainPagePlaylists(userId uint) (models.Playlists, error) {
	var plist models.Playlists

	rows, err := s.conn.Table("kinopoisk.playlists p").
		Select("p.id, p.name, p.public, p.user_id, coalesce(sub.user_id, 0)").
		Joins("inner join kinopoisk.users u1 on u1.id = p.user_id").
		Joins("left join ? sub on sub.playlist_id = p.id",
			s.conn.Table("kinopoisk.subscriptions s").
				Select("s.playlist_id, s.user_id").
				Joins("inner join kinopoisk.users u2 on s.user_id = u2.id").
				Where("u2.id = ?", userId).SubQuery()).
		Where("u1.username = 'admin' and p.public = true").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uid uint
		p := new(models.Playlist)

		err := rows.Scan(&p.Id, &p.Name, &p.Public, &p.UserId, &uid)
		if err != nil {
			return nil, err
		}

		if uid == userId {
			p.IsSubscribed = true
		}

		plist = append(plist, *p)
	}

	return plist, nil
}
