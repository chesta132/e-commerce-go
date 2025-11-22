package service

import (
	"context"
	"product-service/internal/lib/categorylib"
	"product-service/internal/model"
	"product-service/internal/repo"
	"slices"

	"github.com/chesta132/e-commerce-go/shared/squery"
	"github.com/chesta132/e-commerce-go/shared/sslicelib"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
	return s.cr.FindFirst(s.ctx, []squery.Where{{Name: "id", Value: id}})
}

func (s *EchoCategory) UpdateById(id string, update model.Category) error {
	return s.cr.UpdateOne(s.ctx, []squery.Where{{Name: "id", Value: id}}, update)
}

func (s *EchoCategory) DeleteById(id string) error {
	return s.cr.DeleteOne(s.ctx, []squery.Where{{Name: "id", Value: id}})
}

func (s *EchoCategory) FindProductsByCategoryIds(ids []string) ([]model.Product, error) {
	return s.cr.FindProductsByCategoryIds(ids)
}

func (s *EchoCategory) FindCategoriesByProductId(id string) ([]model.Category, error) {
	prod, err := gorm.G[model.Product](s.cr.DB()).Where("id = ?", id).Preload("Categories", nil).Select("id").First(s.ctx)
	return prod.Categories, err
}

func (r *EchoCategory) UpdateCategoryRelations(prodId string, addIds, removeIds []string) (notFoundIds []string, err error) {
	product := model.Product{ID: prodId}

	err = r.cr.DB().Transaction(func(tx *gorm.DB) error {
		if len(removeIds) > 0 || len(addIds) > 0 {
			ids := slices.Concat(removeIds, addIds)
			cats, err := r.cr.FindMany(r.ctx, []squery.Where{{Name: "id", Value: ids, Op: "IN"}})
			if err != nil {
				return err
			}
			catIds := sslicelib.Map(cats, func(index int, item model.Category) string { return item.ID })
			notFoundIds = sslicelib.FilterNotInSlice(ids, catIds)
			if len(cats) == 0 {
				return nil
			}
		}

		build := func(ids []string) (builded []model.Category) {
			for _, id := range ids {
				builded = append(builded, model.Category{ID: id})
			}
			return sslicelib.Filter(builded, func(index int, item model.Category) bool { return !slices.Contains(notFoundIds, item.ID) })
		}

		if len(removeIds) > 0 {
			toRemove := build(removeIds)
			if err := tx.Model(&product).Association("Categories").Delete(&toRemove); err != nil {
				return err
			}
		}

		if len(addIds) > 0 {
			toAdd := build(addIds)
			if err := tx.Model(&product).Association("Categories").Append(&toAdd); err != nil {
				return err
			}
		}

		return nil
	})
	return
}
