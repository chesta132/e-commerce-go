package routes

import (
	"product-service/config"
	"product-service/internal/handler"
	"product-service/internal/repo"
	"product-service/internal/service"

	mw "product-service/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (rt *Route) RegisterPreview(group *echo.Group, prodSvc *service.Product) {
	r := repo.NewPreview(rt.db)
	svc := service.NewPreview(r)
	h := handler.NewPreview(svc, prodSvc)

	group.GET("", h.GetInfo)

	group.GET("/:id", h.GetOne)
	group.PUT("/:id", mw.AdminOnly(h.UpdateOne))
	group.DELETE("/:id", mw.AdminOnly(h.DeleteOne))

	group.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Limit: config.MAX_PREVIEW_SIZE,
	}))
	group.POST("", mw.AdminOnly(h.CreateOne))
}
