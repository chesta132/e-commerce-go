package service

import (
	"context"
	"product-service/internal/lib/categorylib"
	"product-service/internal/model"
	"product-service/internal/repo"

	"github.com/chesta132/e-commerce-go/shared/squery"
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

func (s *EchoCategory) CreateOne(payload model.CreateCategoryPayload) (model.Category, error) {
	category := categorylib.FilterToCreate(payload)
	return *category, s.cr.CreateOne(s.ctx, category)
}

func (s *EchoCategory) FindById(id string) (model.Category, error) {
	return s.cr.FindById(s.ctx, id)
}

func (s *EchoCategory) UpdateById(id string, update model.Category) error {
	return s.cr.UpdateOne(s.ctx, []squery.Where{{Name: "id", Value: id}}, update)
}

func (s *EchoCategory) DeleteById(id string) error {
	return s.cr.DeleteOne(s.ctx, []squery.Where{{Name: "id", Value: id}})
}

func (s *EchoCategory) FindProductByCategoryIds(ids []string) ([]model.Product, error) {
	return s.cr.FindProductByCategoryIds(ids)
}
