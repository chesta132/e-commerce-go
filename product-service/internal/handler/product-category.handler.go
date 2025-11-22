package handler

import (
	"fmt"
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/replylib"
	"product-service/internal/lib/validatorlib"
	"product-service/internal/model"
	"product-service/internal/service"
	"strings"

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

func (h *ProductCategory) GetCategories(c echo.Context) error {
	csvc := h.csvc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	prodId := c.Param("prod-id")

	categories, err := csvc.FindCategoriesByProductId(prodId)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	return rp.Success(categories).OkJSON()
}

func (h *ProductCategory) UpdateCategoryRelations(c echo.Context) error {
	csvc := h.csvc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	prodId := c.Param("prod-id")
	var payload model.UpdateCategoryRelationPayload
	if err := c.Bind(&payload); err != nil {
		return rp.Error(sreplylib.CodeBadRequest, "payload: invalid request body", reply.OptErrorPayload{Details: err.Error()}).FailJSON()
	}
	if err := validatorlib.Validate.Struct(payload); err != nil {
		return errorlib.HandleValidateError(err.(validator.ValidationErrors), rp)
	}

	nfIds, err := csvc.UpdateCategoryRelations(prodId, payload.Add, payload.Remove)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}
	categories, err := csvc.FindCategoriesByProductId(prodId)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	if len(nfIds) > 0 {
		rp.Info(fmt.Sprintf("Category(s) with id [%s] not found", strings.Join(nfIds, ", ")))
	}

	return rp.Success(categories).OkJSON()
}
