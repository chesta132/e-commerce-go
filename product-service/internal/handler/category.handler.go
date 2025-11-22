package handler

import (
	"fmt"
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
	svc  *service.Category
	psvc *service.Product
}

func NewCategory(service *service.Category, psvc *service.Product) *Category {
	return &Category{service, psvc}
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

func (h *Category) UpdateOne(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	id := c.Param("id")
	var update model.Category
	if err := c.Bind(&update); err != nil {
		return rp.Error(sreplylib.CodeBadRequest, "payload: invalid request body", reply.OptErrorPayload{Details: err.Error()}).FailJSON()
	}

	err := svc.UpdateById(id, update)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	cat, err := svc.FindById(id)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	return rp.Success(cat).OkJSON()
}

func (h *Category) DeleteOne(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	psvc := h.psvc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	id := c.Param("id")

	prod, err := psvc.FindByCategoryIds([]string{id})
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}
	if len(prod) > 0 {
		return rp.Error(sreplylib.CodeConflict, fmt.Sprintf("%d product(s) still associated with this category", len(prod))).FailJSON()
	}

	err = svc.DeleteById(id)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	return rp.Success(map[string]string{"id": id}).OkJSON()
}
