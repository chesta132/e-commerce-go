package handler

import (
	"user-service/db/user"
	"user-service/internal/lib/errorlib"
	"user-service/internal/lib/replylib"
	"user-service/internal/service"

	"github.com/chesta132/e-commerce-go/shared/sreplylib"
	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/chesta132/goreply/reply"
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
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	u, ok := c.Get("user").(user.User)
	if !ok {
		return rp.Error(sreplylib.CodeUnauthorized, errorlib.ErrInvalidToken.Error()).FailJSON()
	}
	return rp.Success(u).OkJSON()
}

func (h *User) UpdateOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	u, ok := c.Get("user").(user.User)
	if !ok {
		return rp.Error(sreplylib.CodeUnauthorized, errorlib.ErrInvalidToken.Error()).FailJSON()
	}
	svc := h.us.AttachEcho(c)

	updt := user.User{}
	if err := c.Bind(&updt); err != nil {
		return rp.Error(sreplylib.CodeBadRequest, "payload: invalid request body", reply.OptErrorPayload{Details: err.Error()}).FailJSON()
	}

	u, err := svc.FindByIdAndUpdate(u.ID, user.User{Address: updt.Address, FullName: updt.FullName})
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	return rp.Success(u).OkJSON()
}

func (h *User) RequestToBecomeAdmin(c echo.Context) error {
	id := c.Param("id")
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.us.AttachEcho(c)
	vsvc := h.vs.AttachEcho(c)

	err := vsvc.CreateRequestToBecomeAdmin(id, svc)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	// [IMPORTANT] send notif after notif service developed
	// also this handler not registered yet

	rp.NoContent()
	return nil
}

func (h *User) AcceptToBecomeAdmin(c echo.Context) error {
	id := c.Param("id")
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.us.AttachEcho(c)
	vsvc := h.vs.AttachEcho(c)

	err := vsvc.AcceptToBecomeAdmin(id, svc)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	// [IMPORTANT] send notif after notif service developed
	// also this handler not registered yet

	rp.NoContent()
	return nil
}
