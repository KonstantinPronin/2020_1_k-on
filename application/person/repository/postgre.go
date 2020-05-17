package repository

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/person"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/util"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type PersonDatabase struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewPersonDatabase(conn *gorm.DB, logger *zap.Logger) person.Repository {
	return &PersonDatabase{
		conn:   conn,
		logger: logger,
	}
}

func (rep *PersonDatabase) GetById(id uint) (*models.Person, error) {
	per := new(models.Person)
	err := rep.conn.Table("kinopoisk.persons").Where("id = ?", id).First(per).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewNotFoundError(err.Error())
		}
		return nil, err
	}

	per.Films, err = rep.GetFilms(id)
	if err != nil {
		return nil, err
	}

	per.Series, err = rep.GetSeries(id)
	if err != nil {
		return nil, err
	}

	return per, nil
}

func (rep *PersonDatabase) Add(p *models.Person) (*models.Person, error) {
	return p, rep.conn.Table("kinopoisk.persons").Create(p).Error
}

func (rep *PersonDatabase) Update(p *models.Person) (*models.Person, error) {
	savedPerson := new(models.Person)
	rep.conn.Table("kinopoisk.persons").Where("id = ?", p.Id).First(savedPerson)

	if p.Name != "" {
		savedPerson.Name = p.Name
	}
	if p.Occupation != "" {
		savedPerson.Occupation = p.Occupation
	}
	if p.BirthDate != "" {
		savedPerson.BirthDate = p.BirthDate
	}
	if p.BirthPlace != "" {
		savedPerson.BirthPlace = p.BirthPlace
	}
	if p.Image != "" {
		savedPerson.Image = p.Image
	}

	return savedPerson, rep.conn.Table("kinopoisk.persons").Save(savedPerson).Error
}

func (rep *PersonDatabase) GetFilms(id uint) (models.ListsFilm, error) {
	var films models.ListsFilm

	rows, err := rep.conn.Table("kinopoisk.film_actor fa").
		Select("f.id, f.maingenre, f.russianname, f.image, f.country, f.year, f.agelimit, f.rating").
		Joins("inner join kinopoisk.films f on fa.film_id = f.id").
		Where("fa.person_id = ?", id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	film := new(models.ListFilm)
	for rows.Next() {
		err := rows.Scan(&film.ID, &film.MainGenre, &film.RussianName, &film.Image, &film.Country, &film.Year, &film.AgeLimit, &film.Rating)
		if err != nil {
			return nil, err
		}

		film.Type = "films"
		films = append(films, *film)
	}

	return films, err
}

func (rep *PersonDatabase) GetSeries(id uint) (models.ListSeriesArr, error) {
	var series models.ListSeriesArr

	rows, err := rep.conn.Table("kinopoisk.series_actor sa").
		Select("s.id, s.maingenre, s.russianname, s.image, s.country, s.yearfirst, s.yearlast, s.agelimit, s.rating").
		Joins("inner join kinopoisk.series s on sa.series_id = s.id").
		Where("sa.person_id = ?", id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s := new(models.ListSeries)
	for rows.Next() {
		err := rows.Scan(&s.ID, &s.MainGenre, &s.RussianName, &s.Image, &s.Country, &s.YearFirst, &s.YearLast, &s.AgeLimit, &s.Rating)
		if err != nil {
			return nil, err
		}
		s.Type = "series"
		series = append(series, *s)
	}

	return series, nil
}

func (rep *PersonDatabase) GetActorsForFilm(filmId uint) (models.ListPersonArr, error) {
	var persons models.ListPersonArr

	rows, err := rep.conn.Table("kinopoisk.film_actor f").
		Select("p.id, p.name").
		Joins("inner join kinopoisk.persons p on f.person_id = p.id").
		Where("f.film_id = ?", filmId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	per := new(models.ListPerson)
	for rows.Next() {
		err := rows.Scan(&per.Id, &per.Name)
		if err != nil {
			return nil, err
		}

		persons = append(persons, *per)
	}

	return persons, nil
}

func (rep *PersonDatabase) GetActorsForSeries(seriesId uint) (models.ListPersonArr, error) {
	var persons models.ListPersonArr

	rows, err := rep.conn.Table("kinopoisk.series_actor s").
		Select("p.id, p.name").
		Joins("inner join kinopoisk.persons p on s.person_id = p.id").
		Where("s.series_id = ?", seriesId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	per := new(models.ListPerson)
	for rows.Next() {
		err := rows.Scan(&per.Id, &per.Name)
		if err != nil {
			return nil, err
		}

		persons = append(persons, *per)
	}

	return persons, nil
}

func (rep *PersonDatabase) Search(word string, begin, end int) (models.ListPersonArr, error) {
	var persons models.ListPersonArr
	query := util.PlainTextToQuery(word)

	rows, err := rep.conn.Table("kinopoisk.persons").
		Select("id, name, image").
		Where("textsearchable_index_col @@ to_tsquery('russian', ?)", query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	per := new(models.ListPerson)
	for rows.Next() {
		err := rows.Scan(&per.Id, &per.Name, &per.Image)
		if err != nil {
			return nil, err
		}

		persons = append(persons, *per)
	}

	return persons, nil
}
