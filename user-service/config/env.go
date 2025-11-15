package config

import "os"

var (
	DB_USER              = os.Getenv("DB_USER")
	DB_PASSWORD          = os.Getenv("DB_PASSWORD")
	DB_HOST              = os.Getenv("DB_HOST")
	DB_PORT              = os.Getenv("DB_PORT")
	DB_NAME              = os.Getenv("DB_NAME")
	ACCESS_TOKEN_SECRET  = os.Getenv("ACCESS_TOKEN_SECRET")
	REFRESH_TOKEN_SECRET = os.Getenv("REFRESH_TOKEN_SECRET")
	GO_ENV               = os.Getenv("GO_ENV")
	SERVER_PORT          = os.Getenv("PORT")
	SERVICE              = os.Getenv("SERVICE")
)
