package routes

import (
	"user-service/internal/handler"
	"user-service/internal/middleware"
	"user-service/internal/repo"
	"user-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterAuth(g *echo.Group) {
	ur := repo.NewUser(rt.db)
	rr := repo.NewRevoked(rt.db)
	s := service.NewAuth(ur, rr)
	h := handler.NewAuth(s)
	mw := middleware.NewAuth(s)

	g.POST("/sign-in", h.Signin)
	g.POST("/sign-up", h.Signup)
	g.GET("/user", mw.Protected(h.TokenValid))
}
