package main

import (
	"user-service/config"
	"user-service/db"
	"user-service/internal/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := db.Connect()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	ug := e.Group("/users")
	ag := e.Group("/auth")

	r := routes.New(db)
	r.RegisterUser(ug)
	r.RegisterAuth(ag)

	e.Logger.Info("user-service started at :" + config.SERVER_PORT)
	e.Start(":" + config.SERVER_PORT)
}
