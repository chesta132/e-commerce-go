package service

import (
	"context"
	"time"
	"user-service/db/user"
	"user-service/internal/lib/crypto"
	"user-service/internal/lib/errorlib"
	"user-service/internal/lib/query"
	"user-service/internal/lib/token"
	"user-service/internal/lib/validatorlib"
	"user-service/internal/model"
	"user-service/internal/repo"

	"github.com/labstack/echo/v4"
)

type Auth struct {
	ur *repo.User
	rr *repo.Revoked
}

type EchoAuth struct {
	Auth
	c   echo.Context
	ctx context.Context
}

func NewAuth(repo *repo.User, revokedRepo *repo.Revoked) *Auth {
	return &Auth{repo, revokedRepo}
}

func (s *Auth) AttachEcho(c echo.Context) *EchoAuth {
	return &EchoAuth{Auth: Auth{ur: s.ur, rr: s.rr}, c: c, ctx: c.Request().Context()}
}

func (s *EchoAuth) Signin(payload *model.SigninPayload) (user.User, error) {
	err := validatorlib.Validate.Struct(payload)
	if err != nil {
		return user.User{}, err
	}

	u, err := s.ur.FindFirst(s.ctx, []query.Where{{Name: "email", Value: payload.Email}})
	if err != nil {
		return user.User{}, err
	}
	if !crypto.ComparePassword(u.Password, payload.Password) {
		return user.User{}, errorlib.ErrWrongPassword
	}

	return u, nil
}

func (s *EchoAuth) Signup(payload *model.SignupPayload) (user.User, error) {
	err := validatorlib.Validate.Struct(payload)
	if err != nil {
		return user.User{}, err
	}

	_, err = s.ur.FindFirst(s.ctx, []query.Where{{Name: "email", Value: payload.Email}})
	if err == nil {
		return user.User{}, errorlib.ErrEmailRegistered
	}

	hp := crypto.HashPassword(payload.Password)
	u := user.User{FullName: payload.FullName, Email: payload.Email, Password: hp}

	err = s.ur.CreateOne(s.ctx, &u)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (s *EchoAuth) ValidateAuth(access string, refresh string) (user user.User, newAccess, newRefresh string, error error) {
	if _, err := s.rr.FindByValue(s.ctx, refresh); err == nil { // change to redis later
		return user, "", "", errorlib.ErrTokenExpired
	}

	a, ok := token.ParseAccess(access)
	if !ok {
		r, ok := token.ParseRefresh(refresh)
		if !ok {
			return user, "", "", errorlib.ErrInvalidToken
		}

		user, err := s.ur.FindById(s.ctx, r.ID)
		if err == nil {
			newAccess = token.CreateAccess(user)
		}
		if r.Expires.Before(time.Now()) {
			s.rr.CreateByValue(s.ctx, refresh)
			newRefresh = token.CreateRefresh(user)
		}

		return user, newAccess, newRefresh, err
	}

	user, err := s.ur.FindById(s.ctx, a.ID)
	return user, "", "", err
}
