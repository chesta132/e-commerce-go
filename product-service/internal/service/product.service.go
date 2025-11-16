package service

import (
	"context"
	"product-service/config"
	"product-service/db/product"
	"product-service/internal/repo"

	"github.com/labstack/echo/v4"
)

type Product struct {
	pr *repo.Product
}

type EchoProduct struct {
	Product
	ctx context.Context
	c   echo.Context
}

func NewProduct(repoProduct *repo.Product) *Product {
	return &Product{repoProduct}
}

func (s *Product) AttachEcho(c echo.Context) *EchoProduct {
	return &EchoProduct{Product: *s, ctx: c.Request().Context(), c: c}
}

func (s *EchoProduct) SearchByKeyword(keyword string, offset int) ([]product.Product, error) {
	return s.pr.SearchByKeyword(s.ctx, keyword, offset, config.PAGINATION_LIMIT)
}
