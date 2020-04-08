package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/person"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type Person struct {
	persons person.Repository
	logger  *zap.Logger
}

func NewPerson(p person.Repository, logger *zap.Logger) person.UseCase {
	return &Person{
		persons: p,
		logger:  logger,
	}
}

func (usecase *Person) GetById(id uint) (*models.Person, error) {
	if id == 0 {
		return nil, errors.NewInvalidArgument("wrong id")
	}

	return usecase.persons.GetById(id)
}

func (usecase *Person) Add(p *models.Person) (*models.Person, error) {
	if p.Name == "" {
		return nil, errors.NewInvalidArgument("empty name")
	}
	if !usecase.checkDate(p.BirthDate) {
		return nil, errors.NewInvalidArgument("invalid birth date")
	}

	return usecase.persons.Add(p)
}

func (usecase *Person) Update(p *models.Person) (*models.Person, error) {
	if p.Id == 0 {
		return nil, errors.NewInvalidArgument("wrong id")
	}
	if !usecase.checkDate(p.BirthDate) {
		return nil, errors.NewInvalidArgument("invalid birth date")
	}

	return usecase.persons.Update(p)
}

func (usecase *Person) checkDate(date string) bool {
	if date == "" {
		return true
	}

	_, err := time.Parse("2006-01-02", date)
	if err == nil {
		return true
	}

	return false
}

func (usecase *Person) GetActorsForFilm(filmId uint) (models.ListPersonArr, error) {
	if filmId == 0 {
		return nil, errors.NewInvalidArgument("wrong id")
	}

	return usecase.persons.GetActorsForFilm(filmId)
}

func (usecase *Person) GetActorsForSeries(seriesId uint) (models.ListPersonArr, error) {
	if seriesId == 0 {
		return nil, errors.NewInvalidArgument("wrong id")
	}

	return usecase.persons.GetActorsForSeries(seriesId)
}
