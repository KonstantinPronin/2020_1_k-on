package infrastructure

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func InitDatabase(path string) (db *gorm.DB, err error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var dsn string

	for key, val := range viper.AllSettings() {
		dsn = fmt.Sprintf(`%s%s=%s `, dsn, key, val)
	}

	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.DB().Ping(); err != nil {
		return nil, err
	}
	//db.LogMode(true)
	return db, nil
}
