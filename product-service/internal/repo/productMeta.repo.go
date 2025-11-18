package repo

import (
	"context"
	"product-service/internal/model"

	"gorm.io/gorm"
)

type ProductMeta struct {
	db *gorm.DB
}

func NewProductMeta(db *gorm.DB) *ProductMeta {
	return &ProductMeta{db}
}

func (r *ProductMeta) DB() *gorm.DB {
	return r.db
}

func (r *ProductMeta) CreateOne(ctx context.Context, meta *model.ProductMeta) error {
	return gorm.G[model.ProductMeta](r.db).Create(ctx, meta)
}
