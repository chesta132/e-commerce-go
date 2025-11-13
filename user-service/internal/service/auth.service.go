package service

import (
	"context"
	"user-service/db/user"
	"user-service/internal/lib/crypto"
	"user-service/internal/lib/errorlib"
	"user-service/internal/lib/query"
	"user-service/internal/lib/validate"
	"user-service/internal/repo"

	"github.com/labstack/echo/v4"
)

type Auth struct {
	r *repo.User
}

type EchoAuth struct {
	r   *repo.User
	c   echo.Context
	ctx context.Context
}

func NewAuth(repo *repo.User) *Auth {
	return &Auth{repo}
}

func (r *Auth) AttachEcho(c echo.Context) *EchoAuth {
	return &EchoAuth{r: r.r, c: c, ctx: c.Request().Context()}
}

type SigninPayload struct {
	Email      string `validate:"email,required" json:"email"`
	Password   string `validate:"required" json:"password"`
	RememberMe bool   `json:"remember-me"`
}

func (s *EchoAuth) Signin(payload *SigninPayload) (user.User, error) {
	err := validate.Validate.Struct(payload)
	if err != nil {
		return user.User{}, err
	}

	u, err := s.r.FindFirst(s.ctx, []query.Where{{Name: "email", Value: payload.Email}})
	if err != nil {
		return user.User{}, err
	}
	if !crypto.ComparePassword(u.Password, payload.Password) {
		return user.User{}, errorlib.ErrWrongPassword
	}

	return u, nil
}
