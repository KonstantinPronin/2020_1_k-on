package repository

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/models"
	"github.com/go-park-mail-ru/2020_1_k-on/internal/user"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/errors"
	"github.com/jackc/pgx/pgxpool"
)

type PostgresForUser struct {
	connPool *pgxpool.Pool
}

func NewPostgresForUser(cp *pgxpool.Pool) user.Repository {
	return &PostgresForUser{connPool: cp}
}

func (p PostgresForUser) Add(usr *models.User) (err error) {
	sql := "INSERT INTO KINOPOISK.USERS (USERNAME, PASSWORD, EMAIL, IMAGE) VALUES ($1, $2, $3, $4)"
	return p.exec(sql, usr.Username, usr.Password, usr.Email, usr.ImagePath)
}

func (p PostgresForUser) Update(id int64, usr *models.User) error {
	sql := "UPDATE KINOPOISK.USERS SET USERNAME=$1, PASSWORD=$2, EMAIL=$3, IMAGE=$4 WHERE ID=$5"
	return p.exec(sql, usr.Username, usr.Password, usr.Email, usr.ImagePath, id)
}

func (p PostgresForUser) GetById(id int64) (*models.User, error) {
	sql := "SELECT * FROM KINOPOISK.USERS WHERE ID=$1"
	usr := new(models.User)
	err := p.get(sql, usr, id)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (p PostgresForUser) GetByName(login string) (*models.User, error) {
	sql := "SELECT * FROM KINOPOISK.USERS WHERE USERNAME=$1"
	usr := new(models.User)
	err := p.get(sql, usr, login)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (p PostgresForUser) Contains(login string) (bool, error) {
	_, err := p.GetByName(login)

	if err != nil {
		switch err.(type) {
		case *errors.NotFoundError:
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (p PostgresForUser) exec(sql string, args ...interface{}) error {
	conn, err := p.connPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	res, err := conn.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.NewDbInternalError("user was not created")
	}

	return nil
}

func (p PostgresForUser) get(sql string, usr *models.User, args ...interface{}) error {
	conn, err := p.connPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	res, err := conn.Query(context.Background(), sql, args...)
	if err != nil {
		return err
	}
	defer res.Close()

	if res.Next() {
		err = res.Scan(&usr.Id, &usr.Username, &usr.Password, &usr.Email, &usr.ImagePath)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.NewNotFoundError("user was not found")
}
