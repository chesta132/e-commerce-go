package main

import (
	"product-service/config"
	"product-service/db"
	"product-service/internal/handler"
	"product-service/internal/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := db.Connect()
	e := echo.New()
	router := routes.New(db)

	e.HTTPErrorHandler = handler.DefaultErrorHandler
	e.Use(middleware.Logger(), middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	pg := e.Group("/products")
	cg := e.Group("/categories")
	router.RegisterProduct(pg)
	router.RegisterCategory(cg)

	pg.Any("/*", handler.NotFound)
	cg.Any("/*", handler.NotFound)
	e.Any("/*", handler.NotFound)

	e.Logger.Info("product-service started at :" + config.SERVER_PORT)
	e.Start(":" + config.SERVER_PORT)
}
