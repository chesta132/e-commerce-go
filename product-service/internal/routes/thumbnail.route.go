package routes

import (
	"product-service/config"
	"product-service/internal/handler"
	"product-service/internal/repo"
	"product-service/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (rt *Route) RegisterThumbnail(group *echo.Group, prodSvc *service.Product) {
	r := repo.NewThumbnail(rt.db)
	svc := service.NewThumbnail(r)
	h := handler.NewThumbnail(svc, prodSvc)

	group.GET("", h.GetInfo)

	group.GET("/:id", h.GetOne)
	group.PUT("/:id", h.UpdateOne)
	group.DELETE("/:id", h.DeleteOne)

	group.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Limit: config.MAX_THUMBNAIL_SIZE,
	}))
	group.POST("", h.CreateOne)
}
