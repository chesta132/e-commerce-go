package routes

import (
	"user-service/internal/handler"
	"user-service/internal/middleware"
	"user-service/internal/repo"
	"user-service/internal/service"

	"github.com/labstack/echo/v4"
)

func (rt *Route) RegisterUser(g *echo.Group) {
	vr := repo.NewVerif(rt.db)
	vs := service.NewVerif(vr)

	ur := repo.NewUser(rt.db)
	us := service.NewUser(ur)
	uh := handler.NewUser(us, vs)

	rr := repo.NewRevoked(rt.db)
	as := service.NewAuth(ur, rr)

	mw := middleware.NewAuth(as)

	g.Use(mw.Protected)
	g.GET("", uh.GetOne)
	g.PUT("", uh.UpdateOne)
	g.GET("/admin", middleware.AdminOnly(uh.GetOne))
}
