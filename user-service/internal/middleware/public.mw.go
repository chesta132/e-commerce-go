package middleware

import (
	"user-service/db/user"
	"user-service/internal/lib/errorlib"
	"user-service/internal/lib/replylib"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rp := replylib.Client.New(adapter.AdaptEcho(c))
		user, ok := c.Get("user").(user.User)
		if !ok {
			return rp.Error(replylib.CodeUnauthorized, errorlib.ErrInvalidToken.Error()).FailJSON()
		}

		if user.Role != "admin" {
			return rp.Error(replylib.CodeUnauthorized, errorlib.ErrAdminOnly.Error()).FailJSON()
		}

		return next(c)
	}
}
