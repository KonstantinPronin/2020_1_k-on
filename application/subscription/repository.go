package subscription

type Repository interface {
	Subscribe(pid, userId uint) error
	Unsubscribe(pid, userId uint) error
	Subscriptions(userId uint) ([]uint, error)
}
