package routes

import (
	"user-service/internal/handler"
	"user-service/internal/repo"
	"user-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterAuth(g *echo.Group) {
	r := repo.NewUser(rt.db)
	s := service.NewAuth(r)
	h := handler.NewAuth(s)

	g.POST("/sign-in", h.Signin)
}
