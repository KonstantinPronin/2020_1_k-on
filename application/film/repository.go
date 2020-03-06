package film

import (
	"2020_1_k-on/application/models"
)

//Интерфейсы запросов к бд

type Repository interface {
	Create(film *models.Film) (models.Film, bool)
	GetById(id uint) (*models.Film, bool)
	GetByName(name string) (*models.Film, bool)
}
