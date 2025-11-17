package handler

import (
	"user-service/db/user"
	"user-service/internal/lib/errorlib"
	"user-service/internal/lib/replylib"
	"user-service/internal/lib/token"
	"user-service/internal/model"
	"user-service/internal/service"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/chesta132/goreply/reply"
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
	p := model.SigninPayload{}
	if err := c.Bind(&p); err != nil {
		return rp.Error(replylib.CodeBadRequest, "payload: invalid request body", reply.OptErrorPayload{Details: err.Error()}).FailJSON()
	}

	u, err := s.Signin(&p)

	if err != nil {
		return errorlib.HandleSigninError(err, rp)
	}

	return rp.Success(u).SetCookies(
		token.CreateAccessCookie(u, p.RememberMe),
		token.CreateRefreshCookie(u, p.RememberMe),
	).OkJSON()
}

func (h *Auth) Signup(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	s := h.s.AttachEcho(c)
	p := model.SignupPayload{}
	if err := c.Bind(&p); err != nil {
		return rp.Error(replylib.CodeBadRequest, "payload: invalid request body", reply.OptErrorPayload{Details: err.Error()}).FailJSON()
	}

	u, err := s.Signup(&p)

	if err != nil {
		return errorlib.HandleSignupError(err, rp)
	}

	return rp.Success(u).SetCookies(
		token.CreateAccessCookie(u, p.RememberMe),
		token.CreateRefreshCookie(u, p.RememberMe),
	).OkJSON()
}

func (h *Auth) TokenValid(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	user, ok := c.Get("user").(user.User)
	if !ok {
		return rp.Error(replylib.CodeUnauthorized, errorlib.ErrInvalidToken.Error()).FailJSON()
	}
	return rp.Success(user).OkJSON()
}
