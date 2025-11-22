package routes

import (
	"product-service/internal/handler"
	"product-service/internal/middleware"
	"product-service/internal/repo"
	"product-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterProduct(group *echo.Group) {
	pr := repo.NewProduct(rt.db)
	ps := service.NewProduct(pr)
	ph := handler.NewProduct(ps)

	group.GET("/search", ph.SearchByKeyword)
	group.POST("", middleware.AdminOnly(ph.CreateOne))
	group.PUT("/:id", middleware.AdminOnly(ph.UpdateOne))
	group.DELETE("/:id", middleware.AdminOnly(ph.DeleteOne))

	cgroup := group.Group("/:prod-id/categories")
	rt.RegisterProductCategory(cgroup)
	cgroup.Any("/*", handler.NotFound)

	pgroup := group.Group("/:prod-id/previews")
	rt.RegisterPreview(pgroup, ps)
	pgroup.Any("/*", handler.NotFound)
}
