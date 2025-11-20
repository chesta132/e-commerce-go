package routes

import (
	"product-service/internal/handler"
	"product-service/internal/repo"
	"product-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterProduct(group *echo.Group) {
	pr := repo.NewProduct(rt.db)
	ps := service.NewProduct(pr)
	ph := handler.NewProduct(ps)

	group.GET("/search", ph.SearchByKeyword)
	group.POST("", ph.CreateOne)

	thumbGroup := group.Group("/:prod-id/thumbnails")
	rt.RegisterThumbnail(thumbGroup, ps)
	thumbGroup.Any("/*", handler.NotFound)
}
