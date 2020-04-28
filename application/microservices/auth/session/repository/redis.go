package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/auth/session"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/constants"
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

func (sd *SessionDatabase) Add(sessionId string, userId uint) error {
	err := sd.conn.Set(sessionId, userId, constants.CookieDuration).Err()
	if err != nil {
		sd.logger.Error(err.Error())
	}

	return err
}

func (sd *SessionDatabase) GetUserId(sessionId string) (uint, error) {
	res, err := sd.conn.Get(sessionId).Result()
	if err != nil {
		sd.logger.Error(err.Error())
		return 0, err
	}

	userId, err := strconv.Atoi(res)
	if err != nil {
		sd.logger.Error(err.Error())
		return 0, err
	}

	return uint(userId), nil
}

func (sd *SessionDatabase) Delete(sessionId string) error {
	err := sd.conn.Del(sessionId).Err()
	if err != nil {
		sd.logger.Error(err.Error())
	}

	return err
}
