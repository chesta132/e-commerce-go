package handler

import (
	"errors"
	"user-service/internal/lib/replylib"
	"user-service/internal/service"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	s *service.User
}

func NewUser(service *service.User) *User {
	return &User{service}
}

func (*User) handleErr(c echo.Context, err error) {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rp.Error(replylib.CodeNotFound, err.Error()).FailJSON()
		return
	}
	rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
}

func (h *User) GetOne(c echo.Context) error {
	id := c.Param("id")
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	s := h.s.AttachEcho(c)

	user, err := s.FindById(id)
	if err != nil {
		h.handleErr(c, err)
		return nil
	}

	rp.Success(user).OkJSON()
	return nil
}
