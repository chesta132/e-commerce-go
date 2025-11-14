package db

import (
	"fmt"
	"user-service/config"
	"user-service/db/revoked"
	"user-service/db/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)
	db, err := gorm.Open(postgres.Open(url))

	if err != nil {
		panic("user-service: Failed to connect database")
	}

	db.Exec(`
		DO $$
		BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
						CREATE TYPE user_role AS ENUM ('user', 'admin');
				END IF;
		END
		$$;
	`)
	db.AutoMigrate(&user.User{}, &revoked.Revoked{})

	return db
}
