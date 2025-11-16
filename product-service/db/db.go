package db

import (
	"fmt"
	"product-service/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)
	db, err := gorm.Open(postgres.Open(url))

	if err != nil {
		panic("product-service: Failed to connect database")
	}

	Migrate(db)

	return db
}
