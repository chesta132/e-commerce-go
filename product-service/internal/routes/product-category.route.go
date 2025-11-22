package routes

import (
	"product-service/internal/handler"
	"product-service/internal/middleware"
	"product-service/internal/repo"
	"product-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterProductCategory(group *echo.Group) {
	pr := repo.NewProduct(rt.db)
	ps := service.NewProduct(pr)

	cr := repo.NewCategory(rt.db)
	cs := service.NewCategory(cr)

	h := handler.NewProductCategory(ps, cs)

	group.DELETE("", middleware.AdminOnly(h.DeleteCategories))
}
