package actor

import (
	"2020_1_k-on/application/models"
)

type Repository interface {
	GetByID(ID uint) (models.Actor, bool)
	GetByName(Name string) (models.Actor, bool)
	Create(actor models.Actor) (models.Actor, bool)
	Update(actor models.Actor) (models.Actor, bool)
}
