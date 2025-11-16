package handler

import (
	"fmt"
	"product-service/config"
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/replylib"
	"product-service/internal/lib/userlib"
	"product-service/internal/model"
	"product-service/internal/service"
	"strconv"
	"strings"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Product struct {
	svc *service.Product
}

func NewProduct(service *service.Product) *Product {
	return &Product{service}
}

func (h *Product) SearchByKeyword(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	keyword := strings.TrimSpace(c.QueryParam(config.KEYWORD_QUERY[0]))
	for i, v := range config.KEYWORD_QUERY {
		if keyword != "" {
			break
		}
		keyword = strings.TrimSpace(c.QueryParam(v))
		if keyword == "" && i == len(config.KEYWORD_QUERY)-1 {
			return rp.
				Error(replylib.CodeBadRequest, fmt.Sprintf("payload: %s not found in query", strings.Join(config.KEYWORD_QUERY, " | "))).
				FailJSON()
		}
	}

	products, err := svc.SearchByKeyword(keyword, offset)
	if err != nil {
		return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
	}

	return rp.Success(products).OkJSON()
}

func (h *Product) CreateOne(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))

	payload := model.CreateProductPayload{}
	if err := c.Bind(&payload); err != nil {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}

	user, err := userlib.GetAdminDataWithAuth(c.Cookies())
	if err != nil {
		return rp.Error(replylib.CodeBadGateway, err.Error()).FailJSON()
	}

	product, err := svc.CreateProduct(&payload, user)
	if err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			return errorlib.HandleValidateError(err, rp)
		}
		return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
	}
	return rp.Success(product).CreatedJSON()
}
