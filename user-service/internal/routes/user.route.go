package routes

import (
	"user-service/internal/handler"
	"user-service/internal/repo"
	"user-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterUser(g *echo.Group) {
	r := repo.NewUser(rt.db)
	s := service.NewUser(r)
	h := handler.NewUser(s)

	g.GET("/:id", h.GetOne)
}
