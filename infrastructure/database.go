package infrastructure

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

const dsn = "host=localhost port=5432 user=usr password=123 dbname=films sslmode=disable"

func InitDatabaseConnection() (*gorm.DB, error) {
	config, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Print(err, "\n\n\n\n\n")
		return nil, err
	}
	config.DB()
	err = config.DB().Ping()
	if err != nil {
		fmt.Print(err)
	}

	return config, nil
}
