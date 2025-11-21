package middleware

import (
	"product-service/internal/lib/replylib"
	"product-service/internal/lib/userlib"

	"github.com/chesta132/e-commerce-go/shared/sreplylib"
	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rp := replylib.Client.New(adapter.AdaptEcho(c))
		admin, cookie, err := userlib.GetUserData(c.Cookies(), "/user/admin")
		if err != nil {
			return rp.Error(sreplylib.CodeBadGateway, err.Error()).FailJSON()
		}
		if cookie != "" {
			rp.AddHeader("Set-Cookie", cookie)
		}
		c.Set("user", admin)
		return next(c)
	}
}

func Protected(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rp := replylib.Client.New(adapter.AdaptEcho(c))
		admin, cookie, err := userlib.GetUserData(c.Cookies(), "/user")
		if err != nil {
			return rp.Error(sreplylib.CodeBadGateway, err.Error()).FailJSON()
		}
		if cookie != "" {
			rp.AddHeader("Set-Cookie", cookie)
		}
		c.Set("user", admin)
		return next(c)
	}
}
