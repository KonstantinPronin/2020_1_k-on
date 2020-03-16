package session

type Repository interface {
	Add(sessionId string, userId int64) error
	GetUserId(sessionId string) (userId int64, err error)
	Delete(sessionId string) error
}
