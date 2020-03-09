package repository

import (
	"2020_1_k-on/application/film"
	"2020_1_k-on/application/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
)

//Интерфейсы запросов к бд

type PostgresForFilms struct {
	connPool *pgxpool.Pool
}

func NewPostgresForFilms(cp *pgxpool.Pool) film.Repository {
	return &PostgresForFilms{connPool: cp}
}

func (p PostgresForFilms) Create(film *models.Film) (models.Film, bool) {
	sql := "INSERT INTO FILMS (NAME, AGELIMIT, IMAGE) VALUES ($1, $2, $3)"
	p.exec(sql, film.Name, film.AgeLimit, film.Image)
	f, _ := p.GetByName(film.Name)
	return *f, true
}

func (p PostgresForFilms) GetById(id uint) (*models.Film, bool) {
	sql := "SELECT * FROM FILMS WHERE ID=$1"
	film := new(models.Film)
	p.get(sql, film, id)
	return film, true
}

func (p PostgresForFilms) GetByName(name string) (*models.Film, bool) {
	sql := "SELECT * FROM FILMS WHERE NAME=$1"
	film := new(models.Film)
	p.get(sql, film, name)
	return film, true
}

func (p PostgresForFilms) GetFilmsArr(begin, end uint) (*models.Films, bool) {
	sql := "SELECT * FROM FILMS LIMIT $1 OFFSET $2"
	films := new(models.Films)
	p.getInterval(sql, films, begin, end)
	return films, true
}

func (p PostgresForFilms) exec(sql string, args ...interface{}) (models.Film, error) {
	conn, err := p.connPool.Acquire(context.Background())
	if err != nil {
		return models.Film{}, err
	}
	defer conn.Release()

	res, err := conn.Exec(context.Background(), sql, args...)
	if err != nil {
		fmt.Print(err)
		return models.Film{}, err
	}
	fmt.Print("RES OK")
	if res.RowsAffected() == 0 {
		return models.Film{}, errors.New("film was not created")
	}

	return models.Film{}, nil
}

func (p PostgresForFilms) get(sql string, film *models.Film, args ...interface{}) error {
	conn, err := p.connPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	res, err := conn.Query(context.Background(), sql, args[0])
	if err != nil {
		fmt.Print(err)
		return err
	}
	defer res.Close()
	if res.Next() {
		err = res.Scan(&film.ID, &film.Name, &film.AgeLimit, &film.Image)
		if err != nil {
			fmt.Print(err)
			return err
		}

		return nil
	}

	return errors.New("user was not found")
}

func (p PostgresForFilms) getInterval(sql string, films *models.Films, begin, end uint) error {
	conn, err := p.connPool.Acquire(context.Background())
	if err != nil {
		fmt.Print(err)
		return err
	}
	defer conn.Release()

	res, err := conn.Query(context.Background(), sql, begin, end)
	if err != nil {
		fmt.Print(err)
		return err
	}
	defer res.Close()
	for res.Next() {
		var id uint
		var name string
		var ageLimit int
		var image string
		err := res.Scan(&id, &name, &ageLimit, &image)
		if err != nil {
			return err
		}
		*films = append(*films, models.Film{
			ID:          id,
			Name:        name,
			AgeLimit:    ageLimit,
			Image:       image,
			ImageBase64: "",
		})
	}
	return errors.New("film was not found")
}
