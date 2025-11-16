package routes

import (
	"product-service/internal/handler"
	"product-service/internal/repo"
	"product-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterProduct(productGroup *echo.Group) {
	pr := repo.NewProduct(rt.db)
	ps := service.NewProduct(pr)
	ph := handler.NewProduct(ps)

	productGroup.GET("/search", ph.SearchByKeyword)
}
