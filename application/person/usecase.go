package person

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type UseCase interface {
	GetById(id uint) (*models.Person, error)
	Add(p *models.Person) (*models.Person, error)
	Update(p *models.Person) (*models.Person, error)
	GetActorsForFilm(filmId uint) (models.ListPersonArr, error)
	GetActorsForSeries(seriesId uint) (models.ListPersonArr, error)
	Search(word string, query map[string][]string) (models.ListPersonArr, error)
}
