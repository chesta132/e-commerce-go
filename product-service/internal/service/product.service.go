package service

import (
	"context"
	"product-service/config"
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/imagelib"
	"product-service/internal/lib/productlib"
	"product-service/internal/lib/validatorlib"
	"product-service/internal/model"
	"product-service/internal/repo"

	"github.com/chesta132/e-commerce-go/shared/smodel"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

// This func will limit as pagination limit's config + 1
func (s *EchoProduct) SearchByKeyword(keyword string, offset int, categoryIds []string) (products []model.Product, notFoundCatIds []string, err error) {
	return s.pr.SearchByKeyword(s.ctx, keyword, repo.SearchByKeywordOptions{
		Offset:      offset,
		Limit:       config.PAGINATION_LIMIT + 1,
		CategoryIds: categoryIds,
	})
}

func (s *EchoProduct) CreateProduct(payload *model.CreateProductPayload, admin *smodel.User) (*model.Product, error) {
	if v := validatorlib.Validate.Struct(payload); v != nil {
		return nil, v
	}

	if payload.Image == nil {
		payload.Image = &model.Image{}
	}

	payload.Image.ID = ""
	if payload.Image.Image == nil {
		img, err := imagelib.ReadDefaultImage()
		if err != nil {
			return nil, err
		}
		payload.Image.Image = img
	}
	db := s.pr.DB()

	var data *model.Product
	err := db.Transaction(func(tx *gorm.DB) error {
		cr := repo.NewCategory(tx)
		existingCat, err := cr.FindManyByIds(s.ctx, payload.CategoryIds)
		if err != nil {
			return err
		}
		if len(existingCat) < 1 {
			return errorlib.ErrNoCategoryToCreate
		}

		ir := repo.NewImage(tx)
		if err := ir.CreateOne(s.ctx, payload.Image); err != nil {
			return err
		}

		mr := repo.NewProductMeta(tx)
		meta := productlib.FilterMetaToCreate(payload)
		if err := mr.CreateOne(s.ctx, meta); err != nil {
			return err
		}

		data = productlib.FilterToCreate(payload, *meta, admin.ID, existingCat)
		data.Meta = *meta

		pr := repo.NewProduct(tx)
		return pr.CreateOne(s.ctx, data)
	})
	return data, err
}
