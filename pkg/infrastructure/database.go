package infrastructure

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const dsn = `host=localhost port=5432 user=postgres password=postgres dbname=k_on sslmode=disable`

func InitGorm() (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.DB().Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
