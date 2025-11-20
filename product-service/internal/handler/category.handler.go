package handler

import (
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/replylib"
	"product-service/internal/lib/validatorlib"
	"product-service/internal/model"
	"product-service/internal/service"

	"github.com/chesta132/e-commerce-go/shared/sreplylib"
	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/chesta132/goreply/reply"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Category struct {
	svc *service.Category
}

func NewCategory(service *service.Category) *Category {
	return &Category{service}
}

func (h *Category) CreateOne(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))

	var payload model.CreateCategoryPayload
	if err := c.Bind(&payload); err != nil {
		return rp.Error(sreplylib.CodeBadRequest, "payload: invalid request body", reply.OptErrorPayload{Details: err.Error()}).FailJSON()
	}
	if err := validatorlib.Validate.Struct(payload); err != nil {
		return errorlib.HandleValidateError(err.(validator.ValidationErrors), rp)
	}

	cat, err := svc.CreateOne(payload)
	if err != nil {
		return rp.Error(sreplylib.CodeServerError, err.Error()).FailJSON()
	}

	return rp.Success(cat).CreatedJSON()
}

func (h *Category) GetOne(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	id := c.Param("id")

	cat, err := svc.FindById(id)
	if err != nil {
		return rp.Error(sreplylib.CodeServerError, err.Error()).FailJSON()
	}

	return rp.Success(cat).OkJSON()
}
