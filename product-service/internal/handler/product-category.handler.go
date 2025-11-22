package handler

import (
	"errors"
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

type ProductCategory struct {
	psvc *service.Product
	csvc *service.Category
}

func NewProductCategory(psvc *service.Product, csvc *service.Category) *ProductCategory {
	return &ProductCategory{psvc, csvc}
}

func (h *ProductCategory) DeleteCategories(c echo.Context) error {
	psvc := h.psvc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	prodId := c.Param("prod-id")
	var payload model.DeleteCategoryRelationPayload
	if err := c.Bind(&payload); err != nil {
		return rp.Error(sreplylib.CodeBadRequest, "payload: invalid request body", reply.OptErrorPayload{Details: err.Error()}).FailJSON()
	}
	if err := validatorlib.Validate.Struct(payload); err != nil {
		return errorlib.HandleValidateError(err.(validator.ValidationErrors), rp)
	}

	err := psvc.DeleteCategories(prodId, payload.CatIds)
	if errors.Is(err, errorlib.ErrCantDeleteAllCategories) {
		return rp.Error(sreplylib.CodeConflict, err.Error()).FailJSON()
	}
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	return rp.Success(map[string][]string{"ids": payload.CatIds}).OkJSON()
}
