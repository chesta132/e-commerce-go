package service

import (
	"context"
	"product-service/config"
	"product-service/internal/lib/imagelib"
	"product-service/internal/lib/productlib"
	"product-service/internal/lib/validatorlib"
	"product-service/internal/model"
	"product-service/internal/repo"

	"github.com/chesta132/e-commerce-go/shared"
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

func (s *EchoProduct) SearchByKeyword(keyword string, offset int) ([]model.Product, error) {
	return s.pr.SearchByKeyword(s.ctx, keyword, offset, config.PAGINATION_LIMIT)
}

func (s *EchoProduct) CreateProduct(payload *model.CreateProductPayload, admin *shared.User) (*model.Product, error) {
	if v := validatorlib.Validate.Struct(payload); v != nil {
		return nil, v
	}

	if payload.Image == nil {
		payload.Image = &model.Image{}
	}
	if payload.Categories == nil {
		payload.Categories = &[]model.Category{}
	}

	payload.Image.ID = ""
	if payload.Image.Image == nil {
		img, err := imagelib.ReadDefaultImage()
		if err != nil {
			return nil, err
		}
		payload.Image.Image = img
	}

	var data *model.Product
	err := s.pr.Transaction(func(tx *gorm.DB) error {
		ir := repo.NewImage(tx)
		if err := ir.CreateOne(s.ctx, payload.Image); err != nil {
			return err
		}

		cr := repo.NewCategory(tx)
		if err := cr.CreateMany(s.ctx, payload.Categories); err != nil {
			return err
		}

		mr := repo.NewProductMeta(tx)
		meta := productlib.FilterMetaToCreate(payload)
		if err := mr.CreateOne(s.ctx, meta); err != nil {
			return err
		}

		data = productlib.FilterToCreate(payload, *meta, admin.ID, *payload.Categories)

		pr := repo.NewProduct(tx)
		return pr.CreateOne(s.ctx, data)
	})
	return data, err
}
