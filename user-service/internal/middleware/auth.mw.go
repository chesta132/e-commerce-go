package middleware

import (
	"errors"
	"net/http"
	"user-service/config"
	"user-service/internal/lib/errorlib"
	"user-service/internal/lib/replylib"
	"user-service/internal/lib/token"
	"user-service/internal/service"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	s *service.Auth
}

func NewAuth(service *service.Auth) *Auth {
	return &Auth{service}
}

func (s *Auth) Protected(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		svc := s.s.AttachEcho(c)
		rp := replylib.Client.New(adapter.AdaptEcho(c))

		ac, err := c.Cookie(config.ACCESS_TOKEN_KEY)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				ac = &http.Cookie{}
			} else {
				return rp.Error(replylib.CodeUnauthorized, err.Error()).FailJSON()
			}
		}

		rc, err := c.Cookie(config.REFRESH_TOKEN_KEY)
		if errors.Is(err, http.ErrNoCookie) {
			return rp.Error(replylib.CodeUnauthorized, errorlib.ErrInvalidToken.Error()).FailJSON()
		}
		if err != nil {
			return rp.Error(replylib.CodeUnauthorized, err.Error()).FailJSON()
		}

		user, na, nc, err := svc.ValidateAuth(ac.Value, rc.Value)
		if err != nil {
			return errorlib.HandleValidateAuthError(err, rp)
		}

		if na != "" {
			rp.SetCookies(token.ToCookie(config.ACCESS_TOKEN_KEY, na, token.AccessExpiry))
		}
		if nc != "" {
			rp.SetCookies(token.ToCookie(config.REFRESH_TOKEN_KEY, nc, token.RefreshExpiry))
		}

		c.Set("user", user)

		return next(c)
	}
}
