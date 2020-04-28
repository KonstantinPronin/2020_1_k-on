package session

type Repository interface {
	Add(sessionId string, userId uint) error
	GetUserId(sessionId string) (userId uint, err error)
	Delete(sessionId string) error
}
