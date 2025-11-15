package handler

import (
	"user-service/internal/lib/errorlib"
	"user-service/internal/lib/replylib"
	"user-service/internal/service"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/labstack/echo/v4"
)

type User struct {
	us *service.User
	vs *service.Verification
}

func NewUser(service *service.User, verifService *service.Verification) *User {
	return &User{service, verifService}
}

func (h *User) GetOne(c echo.Context) error {
	id := c.Param("id")
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.us.AttachEcho(c)

	user, err := svc.FindById(id)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	return rp.Success(user).OkJSON()
}

func (h *User) RequestToBecomeAdmin(c echo.Context) error {
	id := c.Param("id")
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.us.AttachEcho(c)
	vsvc := h.vs.AttachEcho(c)

	user, err := svc.FindById(id)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	err = vsvc.CreateRequestToBecomeAdmin(&user)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	// [IMPORTANT] send notif after notif service developed
	// also this handler not registered yet

	return rp.Success(user).Info("Request sent, please wait to be accepted").OkJSON()
}
