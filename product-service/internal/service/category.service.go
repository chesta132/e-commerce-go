package service

import (
	"context"
	"product-service/internal/model"
	"product-service/internal/repo"

	"github.com/labstack/echo/v4"
)

type Category struct {
	cr *repo.Category
}

type EchoCategory struct {
	Category
	ctx context.Context
	c   echo.Context
}

func NewCategory(repo *repo.Category) *Category {
	return &Category{repo}
}

func (s *Category) AttachEcho(c echo.Context) *EchoCategory {
	return &EchoCategory{Category: *s, ctx: c.Request().Context(), c: c}
}

func (s *EchoCategory) CreateMany(categories *[]model.Category) error {
	return s.cr.CreateMany(s.ctx, categories)
}