package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/internal/session"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"sync"
)

type SessionStorage struct {
	mutex    sync.RWMutex
	sessions map[string]int64
}

func NewSessionStorage() session.Repository {
	return &SessionStorage{
		mutex:    sync.RWMutex{},
		sessions: make(map[string]int64),
	}
}

func (s *SessionStorage) Add(sessionId string, userId int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.sessions[sessionId] = userId
	return nil
}

func (s *SessionStorage) GetUserId(sessionId string) (userId int64, err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	userId, ok := s.sessions[sessionId]

	if ok {
		return userId, nil
	}

	return -1, errors.NewNotFoundError("No such session")
}

func (s *SessionStorage) Delete(sessionId string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.sessions, sessionId)
	return nil
}
