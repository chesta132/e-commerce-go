package service

import (
	"context"
	"user-service/db/user"
	"user-service/db/verification"
	"user-service/internal/lib/errorlib"
	"user-service/internal/repo"

	"github.com/labstack/echo/v4"
)

type Verification struct {
	r *repo.Verification
}

type EchoVerification struct {
	Verification
	c   echo.Context
	ctx context.Context
}

func NewVerif(repo *repo.Verification) *Verification {
	return &Verification{repo}
}

func (s *Verification) AttachEcho(c echo.Context) *EchoVerification {
	return &EchoVerification{Verification: *s, c: c, ctx: c.Request().Context()}
}

func (s *EchoVerification) CreateRequestToBecomeAdmin(id string, userService *EchoUser) error {
	user, err := userService.FindById(id)
	if err != nil {
		return err
	}
	if _, err := s.r.FindByValue(s.ctx, user.ID); err == nil {
		return errorlib.ErrAlreadyRequest
	}
	_, err = s.r.CreateOne(s.ctx, user.ID, verification.RequestToBecomeAdmin)
	return err
}

func (s *EchoVerification) AcceptToBecomeAdmin(id string, userService *EchoUser) error {
	_, err := s.r.FindByValue(s.ctx, id)
	if err != nil {
		return err
	}
	return userService.UpdateById(id, user.User{Role: "admin"})
}
