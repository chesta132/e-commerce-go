package routes

import (
	"product-service/internal/handler"
	"product-service/internal/middleware"
	"product-service/internal/repo"
	"product-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterCategory(group *echo.Group) {
	pr := repo.NewProduct(rt.db)
	psvc := service.NewProduct(pr)

	r := repo.NewCategory(rt.db)
	svc := service.NewCategory(r)
	h := handler.NewCategory(svc, psvc)

	group.GET("/:id", h.GetOne)
	group.POST("", middleware.AdminOnly(h.CreateOne))
	group.DELETE("/:id", middleware.AdminOnly(h.DeleteOne))
	group.PUT("/:id", middleware.AdminOnly(h.UpdateOne))
}
