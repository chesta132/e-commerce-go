package routes

import (
	"user-service/internal/handler"
	"user-service/internal/repo"
	"user-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterUser(g *echo.Group) {
	vr := repo.NewVerif(rt.db)
	vs := service.NewVerif(vr)

	r := repo.NewUser(rt.db)
	s := service.NewUser(r)
	h := handler.NewUser(s, vs)

	g.GET("/:id", h.GetOne)
	g.PUT("/:id", h.UpdateOne)
}
