package handler

import (
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

func (h *Auth) Signin(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	s := h.s.AttachEcho(c)
	p := service.SigninPayload{}
	if err := c.Bind(&p); err != nil {
		rp.Error(replylib.CodeBadRequest, "payload: invalid request body").FailJSON()
		return nil
	}

	u, err := s.Signin(&p)

	if err != nil {
		errorlib.HandleSigninError(err, rp)
		return nil
	}

	rp.Success(u).SetCookies(
		token.CreateAccessCookie(u, p.RememberMe),
		token.CreateRefreshCookie(u, p.RememberMe),
	).OkJSON()
	return nil
}
