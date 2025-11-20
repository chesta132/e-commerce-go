package routes

import (
	"product-service/internal/handler"
	"product-service/internal/repo"
	"product-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterThumbnail(group *echo.Group, prodSvc *service.Product) {
	r := repo.NewThumbnail(rt.db)
	svc := service.NewThumbnail(r)
	h := handler.NewThumbnail(svc, prodSvc)

	group.GET("/:id", h.GetOne)
	group.POST("", h.CreateOne)
}
