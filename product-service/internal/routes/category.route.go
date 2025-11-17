package routes

import (
	"product-service/internal/handler"
	"product-service/internal/repo"
	"product-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterCategory(group *echo.Group) {
	r := repo.NewCategory(rt.db)
	svc := service.NewCategory(r)
	h := handler.NewCategory(svc)

	group.GET("/:id", h.GetOne)
	group.POST("", h.CreateOne)
}
