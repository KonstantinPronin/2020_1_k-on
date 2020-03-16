package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/internal/session"
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"strconv"
)

type SessionDatabase struct {
	conn   *redis.Client
	logger *zap.Logger
}

func NewSessionDatabase(conn *redis.Client, logger *zap.Logger) session.Repository {
	return &SessionDatabase{conn: conn, logger: logger}
}

func (sd *SessionDatabase) Add(sessionId string, userId int64) error {
	err := sd.conn.Set(sessionId, userId, session.CookieDuration).Err()
	if err != nil {
		sd.logger.Error(err.Error())
	}

	return err
}

func (sd *SessionDatabase) GetUserId(sessionId string) (int64, error) {
	res, err := sd.conn.Get(sessionId).Result()
	if err != nil {
		sd.logger.Error(err.Error())
		return -1, err
	}

	userId, err := strconv.Atoi(res)
	if err != nil {
		sd.logger.Error(err.Error())
		return -1, err
	}

	return int64(userId), nil
}

func (sd *SessionDatabase) Delete(sessionId string) error {
	err := sd.conn.Del(sessionId).Err()
	if err != nil {
		sd.logger.Error(err.Error())
	}

	return err
}
