package user

type UseCase interface {
	Login(login string, password string) (sessionId string, csrfToken string, err error)
	Check(sessionId string) (userId uint, err error)
	Logout(sessionId string) error
}
