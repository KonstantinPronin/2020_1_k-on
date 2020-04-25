package film

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
)

//Интерфейсы запросов к бд

type Repository interface {
	Create(film *models.Film) (models.Film, bool)
	GetById(id uint) (*models.Film, bool)
	GetByName(name string) (*models.Film, bool)
	GetFilmsArr(begin, end uint) (*models.Films, bool)
	GetFilmGenres(fid uint) (models.Genres, bool)
}
