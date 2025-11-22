package handler

import (
	"fmt"
	"product-service/config"
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/replylib"
	"product-service/internal/model"
	"product-service/internal/service"
	"strconv"
	"strings"

	"github.com/chesta132/e-commerce-go/shared/smodel"
	"github.com/chesta132/e-commerce-go/shared/sreplylib"
	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/chesta132/goreply/reply"
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
	categoryIds := c.QueryParams()["category-id"]
	keyword := strings.TrimSpace(c.QueryParam(config.KEYWORD_QUERY[0]))
	for i, v := range config.KEYWORD_QUERY {
		if keyword != "" {
			break
		}
		keyword = strings.TrimSpace(c.QueryParam(v))
		if keyword == "" && i == len(config.KEYWORD_QUERY)-1 {
			return rp.
				Error(sreplylib.CodeBadRequest, fmt.Sprintf("payload: %s not found in query", strings.Join(config.KEYWORD_QUERY, " | "))).
				FailJSON()
		}
	}

	products, nfCatIds, err := svc.SearchByKeyword(keyword, offset, categoryIds)
	if err != nil {
		return rp.Error(sreplylib.CodeServerError, err.Error()).FailJSON()
	}
	if len(nfCatIds) > 0 {
		rp.Info(fmt.Sprintf("Category(s) with id [%s] not found", strings.Join(nfCatIds, ", ")))
	}

	return rp.Success(products).PaginateCursor(config.PAGINATION_LIMIT, offset).OkJSON()
}

func (h *Product) GetOne(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	id := c.Param("id")

	product, err := svc.FindByIdWithRelation(id, []string{"Meta"})
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	return rp.Success(product).OkJSON()
}

func (h *Product) CreateOne(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))

	payload := model.CreateProductPayload{}
	if err := c.Bind(&payload); err != nil {
		return rp.Error(sreplylib.CodeBadRequest, "payload: invalid request body", reply.OptErrorPayload{Details: err.Error()}).FailJSON()
	}

	user := c.Get("user").(*smodel.User)

	product, err := svc.CreateProduct(&payload, user)
	if err != nil {
		return errorlib.HandleCreateProductError(err, rp)
	}
	return rp.Success(product).CreatedJSON()
}

func (h *Product) UpdateOne(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	id := c.Param("id")

	var update model.Product
	if err := c.Bind(&update); err != nil {
		return rp.Error(sreplylib.CodeBadRequest, "payload: invalid request body", reply.OptErrorPayload{Details: err.Error()}).FailJSON()
	}
	update.SKU = ""

	err := svc.UpdateByid(id, update)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	product, err := svc.FindById(id)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	return rp.Success(product).OkJSON()
}

func (h *Product) DeleteOne(c echo.Context) error {
	svc := h.svc.AttachEcho(c)
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	id := c.Param("id")

	err := svc.DeleteById(id)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	return rp.Success(map[string]string{"id": id}).OkJSON()
}
